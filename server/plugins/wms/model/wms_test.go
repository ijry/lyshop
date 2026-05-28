package model

import "testing"

func TestDocTypeAndStatusConstants(t *testing.T) {
	if !IsValidDocType(DocTypeInbound) {
		t.Fatalf("expected %s to be valid doc type", DocTypeInbound)
	}
	if !IsValidDocType(DocTypeOutbound) {
		t.Fatalf("expected %s to be valid doc type", DocTypeOutbound)
	}
	if IsValidDocType("unknown") {
		t.Fatalf("unexpected doc type accepted")
	}

	if !IsValidDocStatus(DocStatusDraft) {
		t.Fatalf("expected %s to be valid doc status", DocStatusDraft)
	}
	if !IsValidDocStatus(DocStatusCompleted) {
		t.Fatalf("expected %s to be valid doc status", DocStatusCompleted)
	}
	if !IsValidDocStatus(DocStatusCanceled) {
		t.Fatalf("expected %s to be valid doc status", DocStatusCanceled)
	}
	if IsValidDocStatus("invalid") {
		t.Fatalf("unexpected doc status accepted")
	}
}
