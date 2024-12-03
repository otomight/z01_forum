import { extractAttributes } from "./tools/attribute.js";
import { addToElemNumber } from "./tools/math.js";

// match with LikeDislikePostRequestAjax struct in server
interface ReactionRequest {
	post_id:	number;
	user_id:	number;
};

// match with LikeRequest and HTML
interface PostAttributeMap {
	post_id:	string
	user_id:	string
};

// match with LikeDislikePostResponseAjax struct in server
interface ReactionResponse {
	added:		boolean;
	deleted:	boolean;
	replaced:	boolean;
};

const	REACTION_POST_ATTRIBUTE_MAP: PostAttributeMap = {
	post_id:	'post-id',
	user_id:	'current-user-id'
};

function buildRequest(post: HTMLElement): ReactionRequest | null {
	const	data:		PostAttributeMap | null = (
		extractAttributes<PostAttributeMap>(
			post,
			REACTION_POST_ATTRIBUTE_MAP
		)
	);
	let		request:	ReactionRequest;

	if (!data)
		return null;
	request = {
		post_id: parseInt(data.post_id, 10),
		user_id: parseInt(data.user_id, 10)
	};
	return request;
}

async function fetchRequest(
	action:		string,
	request:	ReactionRequest
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
		return null;
	} else if (!response.ok)
		throw new Error(response.status + ' The request failed');
	return response
}

async function sendReactionRequest(
	post:	HTMLElement,
	action:	string
): Promise<ReactionResponse | null> {

	const	request:	ReactionRequest | null = (
		buildRequest(post)
	);
	let		response:	Response | null

	if (request == null)
		return null;
	try {
		response = await fetchRequest(action, request);
		if (!response)
			return null;
		return response.json()
	} catch (error) {
		console.error('Error:', error);
		alert('Something went wrong, please try again.');
		return null
	}
}

function addToButtonValue(button: HTMLButtonElement, nb: number) {
	const	buttonCount:	HTMLElement | null = (
		button.querySelector('.reaction-count') as HTMLElement | null
	);

	if (!buttonCount) {
		console.error("Element with class reaction-count not found");
		return;
	}
	addToElemNumber(buttonCount, nb)
}

async function sendReaction(
	post:			HTMLElement,
	action:			string,
	button:			HTMLButtonElement,
	oppositeButton:	HTMLButtonElement
) {
	const		response:	ReactionResponse | null = (
		await sendReactionRequest(post, action)
	);

	if (response == null)
		return;
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

function handleReactionButton(event: Event) {
	const	target:			HTMLElement | null = (
		event.target as HTMLElement | null
	);
	let		button:			HTMLButtonElement | null;
	let		post:			HTMLElement | null;
	let		oppositeButton:	HTMLButtonElement | null;
	let		action:			string;

	button = target?.closest('button') as HTMLButtonElement | null
	if (!button)
		return;
	if (button.classList.contains('like-button') ||
								button.classList.contains('dislike-button')) {
		post = button.closest('.post');
		if (!post) {
			console.error("ReactionButton: element 'post' not found.");
			return;
		}
		if (button.classList.contains('like-button')) {
			action = "/post/like";
			oppositeButton = post.querySelector('.dislike-button');
		} else {
			action = "/post/dislike";
			oppositeButton = post.querySelector('.like-button');
		}
		if (!oppositeButton) {
			console.error("ReactionButton: element 'opposite-button' not found.");
			return;
		}
		sendReaction(post, action, button, oppositeButton);
	}
}

document.addEventListener('DOMContentLoaded', () => {
	const	buttonsContainer:	HTMLElement | null = (
		document.getElementById('buttons-container')
	);
	if (!buttonsContainer) {
		console.error("No div buttons-container found.");
		return;
	}
	buttonsContainer?.addEventListener('click',
		(event) => handleReactionButton(event))
});
