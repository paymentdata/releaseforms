package main

//ReleaseTemplate is the go html/template used to render the form.
//Can publish HTML, though current workflow is to send the rendered HTML payload
//  to the athenapdf microservice for signing/distribution.
const ReleaseTemplate string = `
<!DOCTYPE html>
<html>
<head>
<title>PDS SOFTWARE RELEASE FORM</title>
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
<script src="https://code.jquery.com/jquery-3.4.1.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>
</head>
<body>
<!-- Form Name -->
<container>
  <div class="row justify-content-start mt-5 text-dark">
    <div class="col">
      <h1 class="text-center text-dark">PDS SOFTWARE RELEASE FORM</h1>
    </div>
  </div>
<!-- Prepended text-->
   <div class="border-bottom">
    <div class="row mt-4 mx-5">
      <div class="col-xl-12">
        <p class="paragraph-txt d-inline"><b>Release Date:</b> {{ .Date.Format "Jan 02, 2006 15:04:05 UTC" }}</p>
      </div>
    </div>
    <div class="row mt-2 mx-5">
      <div class="col-xl-12">
        <!-- Prepended text -->
        <p class="paragraph-txt d-inline"><b>Product:</b> {{ .Product }}</p>
      </div>
    </div>
    </div>

  <div class="row mt-3 mx-5 text-dark">
  <div class="col">
    <!-- Prepended text-->
    <div>
        <div>
      
        {{ range $key, $value := .Changes }}
        <h3 class="sub-heading- my-4">Change Item {{ inc $key }}: {{ $value.Title }}</h4>
        <p class="mb-2"><b>PR:</b> [{{ $value.ID }}]</p>
        <p class="mb-2"><b>Commit:</b> [{{ $value.CommitSHA }}]</p>
        <p class="paragraph-txt"><b>Written By:</b> {{ $value.Developer }}</p>
        <p class="paragraph-txt"><b>Description of Issue:</b> {{ $value.SummaryOfChangesNeeded }}</p>
        <p class="paragraph-txt"><b>Description of Solution:</b> {{ $value.SummaryOfChangesImplemented }}</p>
        {{ if $value.IssueID }}
        <p class="paragraph-txt"><b>Issue ID: </b>{{ $value.IssueID }}</p>
        {{ end }}
        {{ if $value.ApprovedBy }}
        <p class="paragraph-txt"><b>Approved by:</b> {{ $value.ApprovedBy }}</p>
        {{ end }} 
          <hr></hr>
        {{ end }}  
        </div>

    </div>

    </div>
  </div>
  </div>
  <div class="row mx-5 mt-3">
    <div class="col">
    <h3 class="sub-heading- my-3">Additional Information</h3>
      <!-- Prepended text -->
      <div>
        <p class="paragraph-txt d-inline my-2"><b>Back-Out Procedure:</b> {{ .BackOutProc }}</p>
      </div>

      <!-- Prepended text -->
      <div>
          <p class="paragraph-txt d-inline my-2">
          <b>Impacts PCI Compliance:</b> {{ .PCIImpact }}
          </p>
      </div>

      <!-- Prepended text -->
      <div>
        <p class="paragraph-txt d-inline my-2"><b>OWASP Impact:</b> {{ .OWASPImpact }}</p>
      </div>
    </div>
  </div>
</container>
</body>
</html>
`
