{
  "version": 2,
  "builds": [
    {
      "src": "web/*.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    { "src": "/", "dest": "/web/index.go" }
  ]
}
