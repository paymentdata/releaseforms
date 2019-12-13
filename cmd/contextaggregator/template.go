package main

//ReleaseTemplate is the go html/template used to render the form.
//Can publish HTML, though current workflow is to send the rendered HTML payload
//  to the athenapdf microservice for signing/distribution.
const ReleaseTemplate string = `<form class="form-horizontal">
<fieldset>

<!-- Form Name -->
<legend>PDS SOFTWARE RELEASE FORM</legend>

<!-- Prepended text-->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Release Date:</label>
  <div class="col-md-4">
    <div class="input-group">
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .Date }}" type="text">
    </div>
  </div>
</div>

<!-- Prepended text -->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Product:</label>
  <div class="col-md-4">
    <div class="input-group">
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .Product }}" type="text">
    </div>
  </div>
</div>

<!-- Prepended text-->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext"><h1>Included Changes:</h1></label>
  <div class="col-md-4">
    <div class="input-group">
    {{ range $key, $value := .Changes }}
    <h2>Change Item {{ inc $key }}: {{ $value.Title }}</h2><br><b>PR:[{{ $value.ID }}] Commit:[{{ $value.CommitSHA }}]</b></li>
    <h3>Written By:</h3>
    <h5>{{ $value.Developer }}</h5>
    <h3>Summary of changes:</h3>
    <h5>Description of Issue:</h5>
    <p>{{ $value.SummaryOfChangesNeeded }}</p>
    <h5>Description of Solution:</h5>
    <p>{{ $value.SummaryOfChangesImplemented }}</p>
    {{ if $value.IssueID }}
    <b>Issue ID:{{ $value.IssueID }}</b>
    {{ end }}
    {{ if $value.ApprovedBy }}
    <b>Approved by: {{ $value.ApprovedBy }}</b>
    {{ end }}
    <hr><hr>
  {{ end }}    </div>
  </div>
</div>

<!-- Prepended text -->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Back-Out Procedure:</label>
  <div class="col-md-4">
    <div class="input-group">
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .BackOutProc }}" type="text">
    </div>
  </div>
</div>

<!-- Prepended text -->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Impacts PCI Compliance:</label>
  <div class="col-md-4">
    <div class="input-group">
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .PCIImpact }}" type="text">
    </div>
  </div>
</div>

<!-- Prepended text -->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">OWASP Impact:</label>
  <div class="col-md-4">
    <div class="input-group">
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .OWASPImpact }}" type="text">
    </div>
  </div>
</div>

</fieldset>
</form>
`
