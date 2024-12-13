function	handleShowCommentButton(event: Event) {
	const	target:				HTMLElement | null = (
		event.target as HTMLElement | null
	);
	let	button:				HTMLButtonElement | null;
	let	post:				HTMLElement | null;
	let	comments:			HTMLElement | null;

	button = target?.closest('button') as HTMLButtonElement | null
	if (!button)
		return;
	post = button.closest('.post');
	comments = post?.querySelector('.comments') as HTMLElement;
	if (!comments)
		return;
	if (comments.style.display === 'none' || comments.style.display === '') {
		comments.style.display = 'block';
		button.textContent = "Hide comments"
	} else {
		comments.style.display = 'none';
		button.textContent = "Show comments"
	}
}

document.addEventListener('DOMContentLoaded', () => {
	const	target:	HTMLElement | null = (
		document.querySelector('.show-comments')
	)

	target?.addEventListener('click', handleShowCommentButton);
});
