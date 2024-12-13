import { addToElemNumber } from "./tools/math.js";

// match with LikeDislikePostRequestAjax struct in server
interface ReactionRequest {
	elem_id:	number;
	user_id:	number;
};

// match with LikeDislikePostResponseAjax struct in server
interface ReactionResponse {
	added:		boolean;
	deleted:	boolean;
	replaced:	boolean;
};

// match with LikeRequest and HTML
interface ReactionDataSet {
	type:			'post' | 'comment';
	id:				string;
	currentUserId:	string;
};

async function fetchRequest(
	action:		string,
	request:	ReactionRequest
): Promise<Response | null> {
	let		response:	Response;

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
	return response;
}

async function sendReactionRequest(
	dataset:	ReactionDataSet,
	action:		string
): Promise<ReactionResponse | null> {
	const	request:	ReactionRequest = ({
		elem_id: parseInt(dataset.id, 10),
		user_id: parseInt(dataset.currentUserId, 10)
	});
	let		response:	Response | null;

	try {
		response = await fetchRequest(action, request);
		if (!response)
			return null;
		return response.json();
	} catch (error) {
		console.error('Error:', error);
		alert('Something went wrong, please try again.');
		return null;
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
	addToElemNumber(buttonCount, nb);
}

async function sendReaction(
	dataset:		ReactionDataSet,
	action:			string,
	button:			HTMLButtonElement,
	oppositeButton:	HTMLButtonElement
) {
	const		response:	ReactionResponse | null = (
		await sendReactionRequest(dataset, action)
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

function getAction(dataset: ReactionDataSet, liked: boolean): string {
	let	action:	string;

	if (liked)
		action = "/" + dataset.type + "/like";
	else
		action = "/" + dataset.type + "/dislike";
	return action;
}

function handleReactionButton(event: Event) {
	const	target:				HTMLElement | null = (
		event.target as HTMLElement | null
	);
	let		button:				HTMLButtonElement | null;
	let		oppositeButton:		HTMLButtonElement | null;
	let		reactionSection:	HTMLElement | null;
	let		action:				string;
	let		dataset:			ReactionDataSet;

	button = target?.closest('button') as HTMLButtonElement | null;
	if (!button)
		return;
	reactionSection = button.closest('.reaction-section');
	if (!reactionSection)
		return;
	dataset = reactionSection.dataset as unknown as ReactionDataSet;
	if (button.classList.contains('like-button')) {
		action = getAction(dataset, true);
		oppositeButton = reactionSection.querySelector('.dislike-button');
	}
	else {
		action = getAction(dataset, false);
		oppositeButton = reactionSection.querySelector('.like-button');
	}
	if (!oppositeButton)
		return;
	sendReaction(dataset, action, button, oppositeButton);
}

document.addEventListener('DOMContentLoaded', () => {
	const	buttons:	NodeListOf<HTMLElement> = (
		document.querySelectorAll('.like-button, .dislike-button')
	);

	buttons.forEach(button => {
		button.addEventListener('click', handleReactionButton);
	});
});
