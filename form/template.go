package form

const ReleaseTemplate string = `<form class="form-horizontal">
<fieldset>

<!-- Form Name -->
<legend>RELEASE FORM</legend>

<!-- Prepended text-->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Author</label>
  <div class="col-md-4">
    <div class="input-group">
      <span class="input-group-addon">prepend</span>
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .Author }}" type="text">
    </div>
    <p class="help-block">help</p>
  </div>
</div>

<!-- Prepended text-->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Commit</label>
  <div class="col-md-4">
    <div class="input-group">
      <span class="input-group-addon">prepend</span>
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .Commit }}" type="text">
    </div>
    <p class="help-block">help</p>
  </div>
</div>

<!-- Prepended text-->
<div class="form-group">
  <label class="col-md-4 control-label" for="prependedtext">Release Time</label>
  <div class="col-md-4">
    <div class="input-group">
      <span class="input-group-addon">prepend</span>
      <input id="prependedtext" name="prependedtext" class="form-control" placeholder="{{ .Date }}" type="text">
    </div>
    <p class="help-block">help</p>
  </div>
</div>

</fieldset>
</form>
`
