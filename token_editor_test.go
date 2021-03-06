package rewritehtml

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMiss(t *testing.T) {
	output := new(bytes.Buffer)

	editor := NewTokenEditor(output, AfterHead(`<meta name="rewritehtml" content="true" />`))
	editor.Write([]byte(`<html><p></p><pre>`))
	editor.Write([]byte(`moose</pre></html>`))
	editor.Close()

	require.Equal(t, `<html><p></p><pre>moose</pre></html>`, output.String())
}

func TestHit(t *testing.T) {
	output := new(bytes.Buffer)

	editor := NewTokenEditor(output, AfterHead(`<meta name="rewritehtml" content="true" />`))
	editor.Write([]byte(`<html><head></head><pre>`))
	editor.Write([]byte(`moose</pre></html>`))
	editor.Close()

	require.Equal(t, `<html><head><meta name="rewritehtml" content="true" /></head><pre>moose</pre></html>`, output.String())
}

func TestShortCircuitWrite(t *testing.T) {
	output := new(bytes.Buffer)

	zeros := strings.Repeat(`0`, 1024)
	script := strings.Repeat(`var moose; `, 512)

	editor := NewTokenEditor(output, AfterHead(`<meta name="rewritehtml" content="true" />`))
	editor.Write([]byte(`<!DOCTYPE html><html><head>`))
	editor.Write([]byte(`<link rel="icon" type="image/png" href="data:image/png;base64,` + zeros + `</link></head><script>` + script + `</script>`))
	editor.Write([]byte(`<script>` + script + `</script>`))
	editor.Close()

	require.Equal(t, `<!DOCTYPE html><html><head><meta name="rewritehtml" content="true" /><link rel="icon" type="image/png" href="data:image/png;base64,`+zeros+`</link></head><script>`+script+`</script><script>`+script+`</script>`, output.String())
}

func TestShortWrite(t *testing.T) {
	output := new(bytes.Buffer)

	editor := NewTokenEditor(output, AfterHead(`<meta name="rewritehtml" content="true" />`))
	editor.Write([]byte(`<he`))
	editor.Write([]byte(`ad></head><pre>`))
	editor.Write([]byte(`moose</pre>`))
	editor.Close()

	require.Equal(t, `<head><meta name="rewritehtml" content="true" /></head><pre>moose</pre>`, output.String())
}

func TestConcat(t *testing.T) {
	output := new(bytes.Buffer)

	zeros := strings.Repeat(`0`, 1024)
	script := strings.Repeat(`var moose; `, 512)

	editor := NewTokenEditor(output, AfterHead(`<meta name="rewritehtml" content="true" />`))
	editor.Write([]byte(`<!DOCTYPE html><html><head><link rel="icon" type="image/png" href="data:image/png;base64,` + zeros))
	editor.Write([]byte(`</link></head><script>` + script + `</script>`))
	editor.Close()

	require.Equal(t, `<!DOCTYPE html><html><head><meta name="rewritehtml" content="true" /><link rel="icon" type="image/png" href="data:image/png;base64,`+zeros+`</link></head><script>`+script+`</script>`, output.String())
}

func TestCDataWrite(t *testing.T) {
	output := new(bytes.Buffer)

	editor := NewTokenEditor(output, AfterHead(`<meta name="rewritehtml" content="true" />`))
	editor.Write([]byte(`<script>`))
	editor.Write([]byte(`javascript {} <head></head>`))
	editor.Write([]byte(`moose</script><head></head>`))
	editor.Close()

	require.Equal(t, `<script>javascript {} <head></head>moose</script><head><meta name="rewritehtml" content="true" /></head>`, output.String())
}
