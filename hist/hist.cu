#include <stdio.h>
#include <stdlib.h>

#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/shm.h>

#include <string.h>

#include "cub/cub.cuh"

int main(int argc, char *argv[]) {
	float lower_level = atof(argv[1]);
	float upper_level = atof(argv[2]);

	if (lower_level > upper_level) {
		printf("FAIL: lower level is bigger than upper level\n");
		exit(1);
	}

	int num_bins = atoi(argv[3]);
	int num_levels = num_bins + 1;
	int num_samples = atoi(argv[4]);

	int shmID_matBuf1 = atoi(argv[5]);
	int shmID_histogram = atoi(argv[6]);
	
	float* pBase_matBuf1;
	int* pBase_histogram;

	if ((pBase_matBuf1 = (float*) shmat(shmID_matBuf1, NULL, 0)) == (float*) -1) {
		printf("FAIL: cannot get SHM\n");
		exit(1);
	}

	if ((pBase_histogram = (int*) shmat(shmID_histogram, NULL, 0)) == (int*) -1) {
		printf("FAIL: cannot get SHM\n");
		exit(1);
	}

	float* d_matBuf1;
	int* d_histogram;

	cudaMalloc(&d_matBuf1, num_samples * sizeof(float));
	cudaMemcpy(d_matBuf1, pBase_matBuf1, num_samples * sizeof(float), cudaMemcpyHostToDevice);

	cudaMalloc(&d_histogram, num_bins * sizeof(int));

	void* d_temp_storage = NULL;
	size_t temp_storage_bytes = 0;

	cub::DeviceHistogram::HistogramEven(d_temp_storage, temp_storage_bytes, d_matBuf1, d_histogram, num_levels, lower_level, upper_level, num_samples);

	cudaMalloc(&d_temp_storage, temp_storage_bytes);

	cub::DeviceHistogram::HistogramEven(d_temp_storage, temp_storage_bytes, d_matBuf1, d_histogram, num_levels, lower_level, upper_level, num_samples);

	cudaMemcpy(pBase_histogram, d_histogram, (num_bins) * sizeof(int), cudaMemcpyDeviceToHost);

	return 0;
}
