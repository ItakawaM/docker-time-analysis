export function debounce<T>(fn: (value: T) => void, delay = 300) {
	let timeout: ReturnType<typeof setTimeout>;

	return (value: T) => {
		clearTimeout(timeout);
		timeout = setTimeout(() => fn(value), delay);
	};
}
