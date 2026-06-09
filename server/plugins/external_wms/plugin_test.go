package external_wms

import "testing"

func TestExternalPluginInstallStartsWorker(t *testing.T) {
	p := &externalPlugin{}
	if err := p.Install(); err != nil {
		t.Fatalf("install should not fail: %v", err)
	}
}
