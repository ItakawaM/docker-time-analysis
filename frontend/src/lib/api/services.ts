import { postFormData, postJSON } from './api';
import type { ComputeRequest, ComputeResponse, UploadResponse } from './models';

export async function postUpload(file: File): Promise<UploadResponse> {
	const formData = new FormData();
	formData.append('file', file);

	return postFormData<UploadResponse>('/api/upload', formData);
}

export async function postCompute(
	data: ComputeRequest,
	signal: AbortSignal
): Promise<ComputeResponse> {
	return postJSON<ComputeResponse>('/api/compute', data, signal);
}
