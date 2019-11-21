package form

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
  <label class="col-md-4 control-label" for="prependedtext">Included Changes:</label>
  <div class="col-md-4">
    <div class="input-group">
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .Commit }}" type="text">
    </div>
  </div>
</div>

<!-- Prepended text -->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Changes Approved By:</label>
  <div class="col-md-4">
    <div class="input-group">
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .CommitterName }}" type="text">
    </div>
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
