import { toast } from 'sonner';

interface RetryOptions {
  retries?: number;
  initialDelay?: number;
  shouldRetry?: (error) => boolean;
}

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * Retries an asynchronous function with exponential backoff and jitter upon failure.
 * @template T The return type of the function being retried.
 * @param {() => Promise<T>} fn The asynchronous function to execute.
 * @param {object} [options={}] Optional configuration for the retry mechanism.
 * @param {number} [options.retries=3] The maximum number of retry attempts.
 * @param {number} [options.initialDelay=1500] The initial delay (in milliseconds).
 * @param {(error: any) => boolean} [options.shouldRetry] A function to determine if a given error is retryable. Defaults to true.
 * @returns {Promise<T>} A Promise that resolves with the successful result of the function.
 * @throws {Error} Throws the last encountered error if all retries fail or if the error is deemed not retryable.
 */
export async function retry<T>(
  fn: () => Promise<T>,
  options: RetryOptions = {}
): Promise<T> {
  const { retries = 3, initialDelay = 1500, shouldRetry = () => true } = options;

  for (let attempt = 0; attempt <= retries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      if (attempt === retries || !shouldRetry(error)) {
        // If we've reached the last attempt or the error is not retryable, throw it.
        throw error;
      }

      // Calculate delay with exponential backoff and jitter
      const exponentialDelay = initialDelay * Math.pow(2, attempt);
      const jitter = exponentialDelay * 0.2 * Math.random();
      const totalDelay = exponentialDelay + jitter;

      const message = `Rate limited. Retrying in ${Math.round(totalDelay / 1000)}s... (Attempt ${attempt + 1}/${retries})`;
      console.log(message);
      toast.info(message);

      await delay(totalDelay);
    }
  }
  // Should not reach this
  throw new Error('Retry logic failed');
}
