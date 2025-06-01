#include <arm_neon.h>
#include <chrono>
#include <iomanip>
#include <iostream>
#include <random>
#include <vector>

const size_t ARRAY_SIZE = 1'000'000'000;

// Scalar addition (traditional approach)
double scalar_add(const std::vector<int> &a, const std::vector<int> &b,
                  std::vector<int> &result) {
  auto start = std::chrono::high_resolution_clock::now();

  for (size_t i = 0; i < ARRAY_SIZE; ++i) {
    result[i] = a[i] + b[i];
  }

  auto end = std::chrono::high_resolution_clock::now();
  auto duration =
      std::chrono::duration_cast<std::chrono::microseconds>(end - start);
  return duration.count() / 1000.0; // milliseconds
}

// SIMD addition using ARM NEON
double simd_add(const std::vector<int> &a, const std::vector<int> &b,
                std::vector<int> &result) {
  auto start = std::chrono::high_resolution_clock::now();

  size_t simd_size = ARRAY_SIZE - (ARRAY_SIZE % 4); // Process in chunks of 4

  // SIMD processing for aligned data
  for (size_t i = 0; i < simd_size; i += 4) {
    // Load 4 integers at a time
    int32x4_t vec1 = vld1q_s32(&a[i]);
    int32x4_t vec2 = vld1q_s32(&b[i]);

    // Add the vectors
    int32x4_t sum = vaddq_s32(vec1, vec2);

    // Store the result
    vst1q_s32(&result[i], sum);
  }

  // Handle remaining elements (if any)
  for (size_t i = simd_size; i < ARRAY_SIZE; ++i) {
    result[i] = a[i] + b[i];
  }

  auto end = std::chrono::high_resolution_clock::now();
  auto duration =
      std::chrono::duration_cast<std::chrono::microseconds>(end - start);
  return duration.count() / 1000.0; // milliseconds
}

// Verify results are correct
bool verify_results(const std::vector<int> &scalar_result,
                    const std::vector<int> &simd_result) {
  for (size_t i = 0; i < ARRAY_SIZE; ++i) {
    if (scalar_result[i] != simd_result[i]) {
      std::cout << "Mismatch at index " << i << ": scalar=" << scalar_result[i]
                << ", simd=" << simd_result[i] << std::endl;
      return false;
    }
  }
  return true;
}

int main() {
  std::cout << "SIMD Demo: Adding " << ARRAY_SIZE << " numbers\n";
  std::cout << "==========================================\n\n";

  // Initialize data
  std::vector<int> data1(ARRAY_SIZE);
  std::vector<int> data2(ARRAY_SIZE);
  std::vector<int> result_scalar(ARRAY_SIZE);
  std::vector<int> result_simd(ARRAY_SIZE);

  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<int> dis(1, 100);

  std::cout << "Initializing arrays with random data...\n";
  for (size_t i = 0; i < ARRAY_SIZE; ++i) {
    data1[i] = dis(gen);
    data2[i] = dis(gen);
  }
  std::cout << "Done!\n\n";

  // Run scalar addition
  std::cout << "Running scalar addition...\n";
  double scalar_time = scalar_add(data1, data2, result_scalar);
  std::cout << "Scalar time: " << std::fixed << std::setprecision(2)
            << scalar_time << " ms\n\n";

  // Run SIMD addition
  std::cout << "Running SIMD addition...\n";
  double simd_time = simd_add(data1, data2, result_simd);
  std::cout << "SIMD time: " << std::fixed << std::setprecision(2) << simd_time
            << " ms\n\n";

  // Verify results
  std::cout << "Verifying results...\n";
  if (verify_results(result_scalar, result_simd)) {
    std::cout << "✓ Results match!\n\n";
  } else {
    std::cout << "✗ Results don't match!\n\n";
    return 1;
  }

  // Show performance improvement
  double speedup = scalar_time / simd_time;
  std::cout << "Performance Results:\n";
  std::cout << "==================\n";
  std::cout << "Scalar:  " << std::setw(8) << scalar_time << " ms\n";
  std::cout << "SIMD:    " << std::setw(8) << simd_time << " ms\n";
  std::cout << "Speedup: " << std::setw(8) << std::setprecision(2) << speedup
            << "x\n";
  std::cout << "Improvement: " << std::setw(5) << std::setprecision(1)
            << ((speedup - 1.0) * 100.0) << "%\n";

  return 0;
}
