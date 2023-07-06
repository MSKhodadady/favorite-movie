type Payload = { username: string, exp: number, iat: number }

type OnLoginFunc = (payload: Payload, token: string) => void;

interface Movie {
  name: string,
  year: string
}

interface SuggestedMovie {
  name: string
  year: string
  hash: string
}