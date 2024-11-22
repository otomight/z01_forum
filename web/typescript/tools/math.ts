export function addToElemValue(elem: HTMLElement, toAdd: number) {
	let		text:	string | null;
	let		nb:		number;

	text = elem.textContent;
	if (!text)
		return
	nb = parseInt(text);
	if (isNaN(nb))
		return
	elem.textContent = String(nb + toAdd);
}
