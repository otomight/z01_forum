function handleImageTooLarge(event: Event) {
	const	target:		HTMLInputElement | null = (
		event.target as HTMLInputElement
	);
	const	maxSize:	number = 20 * 1024 * 1024;
	const	errorMsg:	HTMLElement | null = (
		document.getElementById('upload-image-error')
	);
	let		file:		File;

	if (!errorMsg || !target || !target.files || target.files.length == 0)
		return
	file = target.files[0]
	if (file.size > maxSize) {
		errorMsg.style.display = 'block';
		target.value = ''
	} else {
		errorMsg.style.display = 'none';
	}
}

document.addEventListener('DOMContentLoaded', () => {
	const	file_upload:	HTMLElement | null = (
		document.getElementById('file-upload')
	)

	file_upload?.addEventListener('change', handleImageTooLarge)
})
