function sendLikeDislike(action, postId, userId) {
	fetch(action, { // request LikeDislikePostRequestAjax
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			post_id: parseInt(postId, 10),
			user_id: parseInt(userId, 10)
		})
	})
	.then(response => {
		if (!response.ok)
			throw new Error('Response from server not ok');
		return response.json();
	})
	.then(data => { // fill with response LikeDislikePostResponseAjax
		document.getElementById('like-count').textContent = data.like_count;
		document.getElementById('dislike-count').textContent = data.dislike_count;
	})
	.catch(error => {
		console.error('Error:', error);
		alert('Something went wrong, please try again.');
	})
}

document.getElementById('likeButton').addEventListener('click', function() {
	const	postId = this.getAttribute('data-post-id');
	const	userId = this.getAttribute('data-user-id');

	console.log("coucou")
	sendLikeDislike("/post/like", postId, userId);
});

document.getElementById('dislikeButton').addEventListener('click', function() {
	const	postId = this.getAttribute('data-post-id');
	const	userId = this.getAttribute('data-user-id');

	sendLikeDislike("/post/dislike", postId, userId);
});
