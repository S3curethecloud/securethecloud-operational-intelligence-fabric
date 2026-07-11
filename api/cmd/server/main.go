package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Asset struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Environment string `json:"environment"`
}

type Process struct {
	Name    string `json:"name"`
	Cmdline string `json:"cmdline"`
	User    string `json:"user"`
}

type Risk struct {
	InitialScore int      `json:"initial_score"`
	Signals      []string `json:"signals"`
}

type RuntimeEventInput struct {
	ID        string                 `json:"id"`
	Source    string                 `json:"source"`
	EventType string                 `json:"event_type"`
	Severity  string                 `json:"severity"`
	Asset     Asset                  `json:"asset"`
	Message   string                 `json:"message"`
	Process   Process                `json:"process"`
	Risk      Risk                   `json:"risk"`
	Raw       map[string]interface{} `json:"-"`
}

type RuntimeEvent struct {
	ID        string                 `json:"id"`
	Source    string                 `json:"source"`
	EventType string                 `json:"event_type"`
	Severity  string                 `json:"severity"`
	Asset     Asset                  `json:"asset"`
	Message   string                 `json:"message"`
	Process   Process                `json:"process"`
	Risk      Risk                   `json:"risk"`
	Raw       map[string]interface{} `json:"raw_event"`
	CreatedAt time.Time              `json:"created_at"`
}

type PolicyDecision struct {
	ID             string                 `json:"id"`
	RuntimeEventID string                 `json:"runtime_event_id"`
	PolicyName     string                 `json:"policy_name"`
	Allow          bool                   `json:"allow"`
	Severity       string                 `json:"severity"`
	Reason         string                 `json:"reason"`
	OPAResult      map[string]interface{} `json:"opa_result"`
	CreatedAt      time.Time              `json:"created_at"`
}

type Incident struct {
	ID             string            `json:"id"`
	Title          string            `json:"title"`
	Severity       string            `json:"severity"`
	Status         string            `json:"status"`
	RiskScore      int               `json:"risk_score"`
	RuntimeEventID string            `json:"runtime_event_id"`
	PolicyDecision PolicyDecision    `json:"policy_decision"`
	Investigation  *InvestigationRun `json:"investigation,omitempty"`
	Approvals      []Approval        `json:"approvals"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type InvestigationRun struct {
	ID                    string    `json:"id"`
	IncidentID            string    `json:"incident_id"`
	Status                string    `json:"status"`
	Summary               string    `json:"summary"`
	Hypothesis            string    `json:"hypothesis"`
	RecommendedAction     string    `json:"recommended_action"`
	Confidence            string    `json:"confidence"`
	EvidenceRefs          []string  `json:"evidence_refs"`
	RequiresHumanApproval bool      `json:"requires_human_approval"`
	CreatedAt             time.Time `json:"created_at"`
}

type Approval struct {
	ID             string    `json:"id"`
	IncidentID     string    `json:"incident_id"`
	Recommendation string    `json:"recommendation"`
	Decision       string    `json:"decision"`
	Reviewer       string    `json:"reviewer"`
	Rationale      string    `json:"rationale"`
	CreatedAt      time.Time `json:"created_at"`
}

type AuditEntry struct {
	ID         string                 `json:"id"`
	Actor      string                 `json:"actor"`
	Action     string                 `json:"action"`
	TargetType string                 `json:"target_type"`
	TargetID   string                 `json:"target_id"`
	Metadata   map[string]interface{} `json:"metadata"`
	CreatedAt  time.Time              `json:"created_at"`
}

type Store struct {
	mu        sync.Mutex
	events    map[string]RuntimeEvent
	policies  map[string]PolicyDecision
	incidents map[string]Incident
	audit     []AuditEntry
}

var store = Store{
	events:    map[string]RuntimeEvent{},
	policies:  map[string]PolicyDecision{},
	incidents: map[string]Incident{},
	audit:     []AuditEntry{},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/v1/events", handleEvents)
	mux.HandleFunc("/v1/incidents", handleIncidents)
	mux.HandleFunc("/v1/incidents/", handleIncidentSubroutes)
	mux.HandleFunc("/v1/audit-log", handleAuditLog)

	addr := ":" + env("PORT", "8080")
	log.Printf("SecureTheCloud OIF API listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createEvent(w, r)
	case http.MethodGet:
		listEvents(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	var input RuntimeEventInput
	if err := json.Unmarshal(body, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid runtime event"})
		return
	}

	id := input.ID
	if id == "" {
		id = newID("evt")
	}

	event := RuntimeEvent{
		ID:        id,
		Source:    input.Source,
		EventType: input.EventType,
		Severity:  input.Severity,
		Asset:     input.Asset,
		Message:   input.Message,
		Process:   input.Process,
		Risk:      input.Risk,
		Raw:       raw,
		CreatedAt: time.Now().UTC(),
	}

	policy := evaluatePolicy(event)
	riskScore := scoreRisk(event, policy)
	incident := Incident{
		ID:             newID("inc"),
		Title:          fmt.Sprintf("%s on %s", title(event.EventType), event.Asset.Name),
		Severity:       maxSeverity(event.Severity, policy.Severity),
		Status:         "awaiting_human_approval",
		RiskScore:      riskScore,
		RuntimeEventID: event.ID,
		PolicyDecision: policy,
		Approvals:      []Approval{},
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	store.mu.Lock()
	store.events[event.ID] = event
	store.policies[policy.ID] = policy
	if riskScore >= 70 {
		store.incidents[incident.ID] = incident
		store.audit = append(store.audit, AuditEntry{
			ID:         newID("aud"),
			Actor:      "system",
			Action:     "incident_created",
			TargetType: "incident",
			TargetID:   incident.ID,
			Metadata: map[string]interface{}{
				"runtime_event_id": event.ID,
				"risk_score":       riskScore,
				"policy_reason":    policy.Reason,
			},
			CreatedAt: time.Now().UTC(),
		})
	}
	store.mu.Unlock()

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"event":           event,
		"policy_decision": policy,
		"risk_score":      riskScore,
		"incident":        incident,
	})
}

func listEvents(w http.ResponseWriter, r *http.Request) {
	store.mu.Lock()
	defer store.mu.Unlock()
	items := make([]RuntimeEvent, 0, len(store.events))
	for _, e := range store.events {
		items = append(items, e)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	writeJSON(w, http.StatusOK, items)
}

func handleIncidents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	items := make([]Incident, 0, len(store.incidents))
	for _, i := range store.incidents {
		items = append(items, i)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	writeJSON(w, http.StatusOK, items)
}

func handleIncidentSubroutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/incidents/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "missing incident id"})
		return
	}
	incidentID := parts[0]

	if len(parts) == 1 && r.Method == http.MethodGet {
		getIncident(w, r, incidentID)
		return
	}

	if len(parts) == 2 && parts[1] == "investigate" && r.Method == http.MethodPost {
		investigateIncident(w, r, incidentID)
		return
	}

	if len(parts) == 2 && parts[1] == "approvals" && r.Method == http.MethodPost {
		approveIncident(w, r, incidentID)
		return
	}

	if len(parts) == 2 && parts[1] == "evidence" && r.Method == http.MethodGet {
		replayEvidence(w, r, incidentID)
		return
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "route not found"})
}

func getIncident(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	incident, ok := store.incidents[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "incident not found"})
		return
	}
	writeJSON(w, http.StatusOK, incident)
}

func investigateIncident(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.Lock()
	incident, ok := store.incidents[id]
	event := store.events[incident.RuntimeEventID]
	store.mu.Unlock()
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "incident not found"})
		return
	}

	payload := map[string]interface{}{
		"incident":        incident,
		"runtime_events":  []RuntimeEvent{event},
		"policy_decision": incident.PolicyDecision,
		"boundary":        "AI summarizes and recommends only; human approval required for high-risk recommendations.",
	}

	investigation, err := callAIService(payload)
	if err != nil {
		investigation = InvestigationRun{
			ID:                    newID("inv"),
			IncidentID:            id,
			Status:                "completed_with_local_fallback",
			Summary:               "Suspicious shell execution was detected in the payment-api workload. OPA policy context indicates this requires human review.",
			Hypothesis:            "Likely unauthorized shell execution or misconfigured operational command inside a sensitive workload.",
			RecommendedAction:     "Request human approval to isolate the workload in a lab environment and collect additional forensic evidence.",
			Confidence:            "medium",
			EvidenceRefs:          []string{event.ID, incident.PolicyDecision.ID},
			RequiresHumanApproval: true,
			CreatedAt:             time.Now().UTC(),
		}
	}

	store.mu.Lock()
	incident.Investigation = &investigation
	incident.UpdatedAt = time.Now().UTC()
	store.incidents[id] = incident
	store.audit = append(store.audit, AuditEntry{
		ID:         newID("aud"),
		Actor:      "ai-service",
		Action:     "investigation_completed",
		TargetType: "incident",
		TargetID:   id,
		Metadata: map[string]interface{}{
			"requires_human_approval": investigation.RequiresHumanApproval,
			"confidence":              investigation.Confidence,
		},
		CreatedAt: time.Now().UTC(),
	})
	store.mu.Unlock()

	writeJSON(w, http.StatusOK, investigation)
}

func approveIncident(w http.ResponseWriter, r *http.Request, id string) {
	var input struct {
		Decision  string `json:"decision"`
		Reviewer  string `json:"reviewer"`
		Rationale string `json:"rationale"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid approval payload"})
		return
	}
	if input.Decision != "approved" && input.Decision != "rejected" && input.Decision != "needs_more_evidence" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "decision must be approved, rejected, or needs_more_evidence"})
		return
	}
	if input.Reviewer == "" || input.Rationale == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "reviewer and rationale are required"})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	incident, ok := store.incidents[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "incident not found"})
		return
	}

	recommendation := "No AI investigation has been recorded yet."
	if incident.Investigation != nil {
		recommendation = incident.Investigation.RecommendedAction
	}

	approval := Approval{
		ID:             newID("app"),
		IncidentID:     id,
		Recommendation: recommendation,
		Decision:       input.Decision,
		Reviewer:       input.Reviewer,
		Rationale:      input.Rationale,
		CreatedAt:      time.Now().UTC(),
	}
	incident.Approvals = append(incident.Approvals, approval)
	incident.Status = "approval_" + input.Decision
	incident.UpdatedAt = time.Now().UTC()
	store.incidents[id] = incident
	store.audit = append(store.audit, AuditEntry{
		ID:         newID("aud"),
		Actor:      input.Reviewer,
		Action:     "human_approval_recorded",
		TargetType: "incident",
		TargetID:   id,
		Metadata: map[string]interface{}{
			"decision":  input.Decision,
			"rationale": input.Rationale,
		},
		CreatedAt: time.Now().UTC(),
	})

	writeJSON(w, http.StatusCreated, approval)
}

func replayEvidence(w http.ResponseWriter, r *http.Request, id string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	incident, ok := store.incidents[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "incident not found"})
		return
	}
	event := store.events[incident.RuntimeEventID]
	chain := []map[string]interface{}{
		{"step": 1, "type": "runtime_event", "source": event.Source, "id": event.ID, "summary": event.Message},
		{"step": 2, "type": "policy_decision", "source": "opa", "id": incident.PolicyDecision.ID, "summary": incident.PolicyDecision.Reason},
		{"step": 3, "type": "risk_score", "source": "api", "summary": fmt.Sprintf("Risk score calculated as %d", incident.RiskScore)},
	}
	if incident.Investigation != nil {
		chain = append(chain, map[string]interface{}{"step": 4, "type": "ai_investigation", "source": "ai-service", "id": incident.Investigation.ID, "summary": incident.Investigation.Summary})
	}
	for _, approval := range incident.Approvals {
		chain = append(chain, map[string]interface{}{"step": len(chain) + 1, "type": "human_approval", "source": approval.Reviewer, "id": approval.ID, "summary": approval.Decision + ": " + approval.Rationale})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"incident_id": id, "chain": chain})
}

func handleAuditLog(w http.ResponseWriter, r *http.Request) {
	store.mu.Lock()
	defer store.mu.Unlock()
	writeJSON(w, http.StatusOK, store.audit)
}

func evaluatePolicy(event RuntimeEvent) PolicyDecision {
	input := map[string]interface{}{"input": event.Raw}
	body, _ := json.Marshal(input)
	opaURL := env("OPA_URL", "http://localhost:8181") + "/v1/data/securethecloud/runtime/decision"
	resp, err := http.Post(opaURL, "application/json", bytes.NewReader(body))
	if err == nil && resp != nil {
		defer resp.Body.Close()
		var decoded map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&decoded); err == nil {
			if result, ok := decoded["result"].(map[string]interface{}); ok {
				allow, _ := result["allow"].(bool)
				severity, _ := result["severity"].(string)
				reason, _ := result["reason"].(string)
				return PolicyDecision{ID: newID("pol"), RuntimeEventID: event.ID, PolicyName: "securethecloud.runtime.decision", Allow: allow, Severity: severity, Reason: reason, OPAResult: decoded, CreatedAt: time.Now().UTC()}
			}
		}
	}

	allow := true
	severity := "low"
	reason := "No policy violation detected"
	if event.EventType == "suspicious_process_exec" && event.Asset.Name == "payment-api" && event.Process.Name == "sh" {
		allow = false
		severity = "high"
		reason = "Unexpected shell execution in sensitive workload"
	}
	if event.Asset.Name == "payment-api" && strings.Contains(event.Process.Cmdline, "secrets") {
		allow = false
		severity = "critical"
		reason = "Secret file access attempt inside payment workload"
	}
	return PolicyDecision{ID: newID("pol"), RuntimeEventID: event.ID, PolicyName: "local_fallback.runtime.decision", Allow: allow, Severity: severity, Reason: reason, OPAResult: map[string]interface{}{"fallback": true}, CreatedAt: time.Now().UTC()}
}

func callAIService(payload map[string]interface{}) (InvestigationRun, error) {
	body, _ := json.Marshal(payload)
	url := env("AI_SERVICE_URL", "http://localhost:8081") + "/v1/investigate"
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return InvestigationRun{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return InvestigationRun{}, fmt.Errorf("ai service returned status %d", resp.StatusCode)
	}
	var investigation InvestigationRun
	if err := json.NewDecoder(resp.Body).Decode(&investigation); err != nil {
		return InvestigationRun{}, err
	}
	if investigation.ID == "" {
		investigation.ID = newID("inv")
	}
	if investigation.CreatedAt.IsZero() {
		investigation.CreatedAt = time.Now().UTC()
	}
	return investigation, nil
}

func scoreRisk(event RuntimeEvent, policy PolicyDecision) int {
	score := event.Risk.InitialScore
	if !policy.Allow {
		score += 10
	}
	if event.Asset.Name == "payment-api" {
		score += 10
	}
	if strings.Contains(event.Process.Cmdline, "secrets") {
		score += 10
	}
	if event.Severity == "high" {
		score += 5
	}
	if event.Severity == "critical" || policy.Severity == "critical" {
		score += 15
	}
	if score > 100 {
		score = 100
	}
	return score
}

func maxSeverity(a, b string) string {
	order := map[string]int{"low": 1, "medium": 2, "high": 3, "critical": 4}
	if order[b] > order[a] {
		return b
	}
	return a
}

func title(s string) string {
	return strings.ReplaceAll(s, "_", " ")
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func newID(prefix string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return prefix + "_" + hex.EncodeToString(buf)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
