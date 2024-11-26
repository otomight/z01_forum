import { extractAttributes } from "./tools/attribute.js";
import { addToElemValue } from "./tools/math.js";

// match with LikeDislikePostRequestAjax struct in server
interface LikeRequest {
	post_id: number;
	user_id: number;
}

// match with LikeRequest and HTML
interface LikeAttributeMap {
	post_id: string
	user_id: string
}

// match with LikeDislikePostResponseAjax struct in server
interface LikeResponse {
	added:		boolean;
	deleted:	boolean;
	replaced:	boolean;
}

const	LIKE_DISLIKE_POST_ATTRIBUTE_MAP: LikeAttributeMap = {
	post_id: 'data-post-id',
	user_id: 'data-user-id'
}

function	buildRequest(button: HTMLElement): LikeRequest | null {
	const	data:		LikeAttributeMap | null = (
		extractAttributes<LikeAttributeMap>(
			button,
			LIKE_DISLIKE_POST_ATTRIBUTE_MAP
		)
	);
	let		request:	LikeRequest;

	if (!data)
		return null
	request = {
		post_id: parseInt(data.post_id, 10),
		user_id: parseInt(data.user_id, 10)
	}
	return request
}

async function	fetchRequest(
	action: string,
	request: LikeRequest
): Promise<Response | null> {
	let		response:	Response

	response = await fetch(action, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(request)
	});
	if (response.status == 401) {
		alert('You must be logged in to use this feature.');
		return null
	} else if (!response.ok)
		throw new Error(response.status + ' The request failed');
	return response
}

async function sendLikeDislikeRequest(
	button: HTMLElement,
	action: string
): Promise<LikeResponse | null> {

	const	request:	LikeRequest | null = (
		buildRequest(button)
	);
	let		response:	Response | null

	if (request == null)
		return null
	try {
		response = await fetchRequest(action, request)
		if (response == null)
			return null
		return response.json()
	} catch (error) {
		console.error('Error:', error);
		alert('Something went wrong, please try again.');
		return null
	}
}

function addToButtonValue(button: HTMLButtonElement, nb: number) {
	const	buttonCount:	HTMLElement | null = (
		button.querySelector('.like-dislike-count') as HTMLElement | null
	);

	if (!buttonCount) {
		console.error("Element with class like-dislike-count not found")
		return
	}
	addToElemValue(buttonCount, nb)
}

async function handleLikeDislikeButton(event: Event, action: string,
									oppositeButton: HTMLButtonElement) {
	const	button:		HTMLButtonElement | null = (
		event.currentTarget as HTMLButtonElement | null
	);
	let		response:	LikeResponse | null

	if (!button)
		return
	response = await sendLikeDislikeRequest(button, action)
	if (response == null)
		return
	if (response.added) {
		addToButtonValue(button, 1);
		button.classList.add('active');
	} else if (response.deleted) {
		addToButtonValue(button, -1);
		button.classList.remove('active');
	}
	if (response.added && response.replaced) {
		addToButtonValue(oppositeButton, -1);
		oppositeButton.classList.remove('active');
	}
}

document.addEventListener('DOMContentLoaded', () => {
	const	likeButton:		HTMLButtonElement = (
		document.getElementById('likeButton') as HTMLButtonElement
	);
	const	dislikeButton:	HTMLButtonElement = (
		document.getElementById('dislikeButton') as HTMLButtonElement
	);

	likeButton.addEventListener('click',
		(event) => handleLikeDislikeButton(event, "/post/like", dislikeButton)
	);
	dislikeButton.addEventListener('click',
		(event) => handleLikeDislikeButton(event, "/post/dislike", likeButton)
	);
});
