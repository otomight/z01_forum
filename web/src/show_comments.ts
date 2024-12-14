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
		button.classList.add('down');
	}
	else {
		comments.style.display = 'none';
		button.classList.remove('down');
	}
}

document.addEventListener('DOMContentLoaded', () => {
	const	targets:	NodeListOf<HTMLElement> | null = (
		document.querySelectorAll('.show-comments')
	)

	targets.forEach((target: HTMLElement | null) => {
		target?.addEventListener('click', handleShowCommentButton);
	})
});
