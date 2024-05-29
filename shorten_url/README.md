# URL Shortening Service

**Requirements**:

**Shorten a URL**: Provide a long URL and get a shortened version.

**Redirect to the Original URL**: When accessing the shortened URL, redirect to the original URL.

**Concurrency**: Ensure the service handles multiple requests concurrently.

**Storage**: Use an in-memory storage (a map) to store the mappings between original and shortened URLs.