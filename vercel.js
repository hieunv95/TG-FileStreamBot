{
    "version": 2,
    "build": {
      "env": {
        "GO_BUILD_FLAGS": "-ldflags '-s -w'"
      }
    },
    "builds": [
      {
        "src": "cmd/fsb/*.go",
        "use": "@vercel/go"
      }
    ],
    "routes": [
      { "src": "/.*", "dest": "/app/fsb" }
    ]
}
  