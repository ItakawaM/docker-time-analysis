export function debounce<T>(fn: (value: T) => void, delay = 300) {
	let timeout: ReturnType<typeof setTimeout>;

	return (value: T) => {
		clearTimeout(timeout);
		timeout = setTimeout(() => fn(value), delay);
	};
}

export function format(n: number, base: number = 0): string {
	return n.toFixed(base);
}
