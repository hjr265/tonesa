{{define "content"}}
	<div class="row">
		<div class="col-sm-7 col-sm-offset-1">
			<div id="uploadContent">
				<div class="imageWrap">
					<div class="throbber">Loading..</div>
				</div>
			</div>
		</div>
		<div class="col-sm-3">
			<section class="messages">
				<div id="messageDraft">
					<form>
						<div class="form-group">
							<input id="inpMessageAuthorName" class="form-control" value="Anonymous">
						</div>
						<div class="form-group">
							<textarea id="inpMessageContent" class="form-control"></textarea>
						</div>
					</form>
				</div>
				<div id="messages"></div>
			</div>
		</div>
	</div>

	<script type="text/template" id="tplMessageList">
		<ul class="list-unstyled"></ul>
	</script>

	<script type="text/template" id="tplMessageListItem">
		<p title="<%- createdAt %>">
			<b><%- authorName %>:</b> <%- content %>
		</p>
	</script>
{{end}}

{{define "scripts"}}
	<script type="text/javascript">
		var upload = {
			id: {{.Upload.ID.Hex}},
			shortID: {{.Upload.ShortID}},
			content: {
				url: {{.Upload.Content.SignedURL}}
			}
		}
	</script>
	<script src="/assets/js/uploadView.js"></script>
{{end}}
