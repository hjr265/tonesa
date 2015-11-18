(function() {
	'use strict'

	var $uploadContentEl = $('#uploadContent')
	var $contentImgWrapEl = $('#uploadContent .imageWrap')
	var contentImgWidth = 0
	var contentImgHeight = 0
	var $contentImgEl = $('<img class="img-responsive">').attr('src', upload.content.url).on('load', function() {
		contentImgWidth = this.width
		contentImgHeight = this.height
		$contentImgWrapEl.html('').append(this)
		$(window).trigger('resize')

		$contentImgEl.on('click', function() {
			var offset = $(this).offset()
			socket.send('touch upload:'+upload.id+' '+(event.pageX - offset.left) / $contentImgEl.width()+' '+(event.pageY - offset.top) / $contentImgEl.height())
		})
	})
	var $sectionMessagesEl = $('section.messages')
	$(window).on('resize', function() {
		$uploadContentEl.height($(window).height()*0.8)
		$sectionMessagesEl.height($uploadContentEl.height())
		$('#messages').height($('section.messages').height() - $messageDraftEl.outerHeight() - 24)

		$contentImgWrapEl.width(Math.min(contentImgWidth * Math.min(0.9, (($uploadContentEl.height() * 0.9) / contentImgHeight)), contentImgWidth))
		$contentImgWrapEl.css({
			marginTop: -$contentImgEl.height()/2,
			marginLeft: -$contentImgEl.width()/2
		})
	})

	var Message = Backbone.Model.extend({})
	var Messages = Backbone.Collection.extend({
		model: Message,
		url: '/api/uploads/'+upload.id+'/messages',
		comparator: function(a, b) {
			switch(true) {
				case a.get('createdAt') > b.get('createdAt'):
					return -1
				case a.get('createdAt') < b.get('createdAt'):
					return 1
				default:
					return 0
			}
		}
	})

	var messages = new Messages()

	messages.fetch()

	var MessageDraftView = Backbone.View.extend({
		events: {
			'keyup textarea': function(event) {
				switch(event.which) {
					case 13:
						this.createMessage()
						break
				}
			}
		},

		initialize: function() {
			this.delegateEvents()
		},

		createMessage: function() {
			var $messageAuthorNameEl = this.$('#inpMessageAuthorName')
			var $messageContentEl = this.$('#inpMessageContent')
			messages.create({
				authorName: $messageAuthorNameEl.val(),
				content: $messageContentEl.val(),
				createdAt: ''
			}, {
				at: 0
			})
			$messageContentEl.val('')
		}
	})

	var MessageList = Backbone.View.extend({
		template: _.template($('#tplMessageList').detach().html()),

		initialize: function() {
			this.listenTo(this.collection, 'add remove sync', this.render)
		},

		render: function() {
			this.$el.html(this.template())
			this.collection.each(function(msg) {
				var messageListItem = new MessageListItem({
					model: msg
				})
				messageListItem.render()
				this.$('ul').append(messageListItem.el)
			}, this)
		}
	})

	var MessageListItem = Backbone.View.extend({
		tagName: 'li',
		template: _.template($('#tplMessageListItem').detach().html()),

		initialize: function() {
			this.listenTo(this.model, 'change', this.render)
		},

		render: function() {
			this.$el.html(this.template(this.model.toJSON()))
		}
	})

	var $messageDraftEl = $('#messageDraft')
	var messageDraftView = new MessageDraftView({
		el: $messageDraftEl[0]
	})
	messageDraftView.render()

	var $messagesEl = $('#messages')
	var messageList = new MessageList({
		el: $messagesEl[0],
		collection: messages
	})
	messageList.render()

	var socket = glue()
	socket.onMessage(function(data) {
		data = data.split(':')
		switch(data[0]) {
			case 'message':
				messages.fetch({
					data: {
						since: _.first(messages.pluck('createdAt')) || ''
					},
					add: true,
					remove: false
				})
				break

			case 'touch':
				var coords = data[1].split(',')
				showTouchBubble(coords)
				break
		}
	})
	socket.send('watch upload:'+upload.id)

	$(window).trigger('resize')

	function showTouchBubble(coords) {
		var $touchSpotEl = $('<div class="touch"></div>').css({
			top: (Math.round(parseFloat(coords[1])*10000)/100)+'%',
			left: (Math.round(parseFloat(coords[0])*10000)/100)+'%'
		})
		$contentImgWrapEl.append($touchSpotEl)
		$touchSpotEl.animate({
			width: 100,
			height: 100,
			marginTop: -52,
			marginLeft: -52,
			opacity: 0
		}, 2000, function() {
			$touchSpotEl.detach()
		})
	}

})()
