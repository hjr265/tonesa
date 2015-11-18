{{define "content"}}
	<div class="row">
		<div class="col-md-12">
			<div class="text-center">
				<h4>Drop an image here to share it</h4>
				<p>Less than 1MB in size</p>
			</div>
		</div>
	</div>
{{end}}

{{define "scripts"}}
	<script src="//cdnjs.cloudflare.com/ajax/libs/dropzone/4.2.0/min/dropzone.min.js"></script>
	<script type="text/javascript">
		(function() {
			'use script'

			new Dropzone(document.body, {
				url: '/api/uploads',
				maxFiles: 1,
				maxFilesize: 1,
				clickable: false,
				init: function() {
					this.on('success', function(file, resp) {
						location = '/u/'+resp.shortID
					})
					this.on('uploadprogress', function(file, progress) {
						$('h4').text('Uploading... ('+progress+'%)')
						$('p').detach()
					})
					this.on('maxfilesexceeded', function() {
						alert('Image file size is too big')
					})
				}
			})
		})()
	</script>
{{end}}
