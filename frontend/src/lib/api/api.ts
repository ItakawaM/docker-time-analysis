export async function postJSON<T>(url: string, data: any, signal?: AbortSignal): Promise<T> {
	const response = await fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(data),
		signal
	});

	if (!response.ok) {
		throw new Error(await response.text());
	}

	return response.json();
}

export async function postFormData<T>(url: string, formData: FormData): Promise<T> {
	const response = await fetch(url, {
		method: 'POST',
		body: formData
	});

	if (!response.ok) {
		throw new Error(await response.text());
	}

	return response.json();
}
