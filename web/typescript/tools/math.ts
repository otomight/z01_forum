export function addToElemValue(elem: HTMLElement, toAdd: number) {
	const	text:	string | null = elem.textContent;
	let		nb:		number;

	if (!text)
		return
	nb = parseInt(text);
	if (isNaN(nb))
		return
	elem.textContent = String(nb + toAdd);
}
