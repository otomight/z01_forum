export function extractAttributes<T>(
	element:		HTMLElement,
	attributeMap:	Record<keyof T, string>
): T | null {
	const	result:			Partial<T> = {};
	let		key:			keyof T;
	let		attributeName:	string;
	let		value:			string | null;

	for (key in attributeMap) {
		attributeName = attributeMap[key];
		value = element.getAttribute(attributeName);
		if (value === null) {
			console.error(`Attribute ${attributeName} not found on element.`);
			return null;
		}
		result[key] = value as T[keyof T];
	}
	return result as T;
}
