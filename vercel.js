{
    "version": 2,
    "builds": [
      {
        "src": "build.sh",
        "use": "@vercel/go",
        "config": { "buildCommand": "sh build.sh" }
      }
    ],
    "routes": [
      { "src": "/.*", "dest": "/app/fsb" }
    ]
}
  