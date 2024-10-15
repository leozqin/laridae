# Laridae
## About
Laridae is a registry and proxy server for agentic tool-calling. The name comes from the family name for Gulls, which are known to use tools.

Laridae is designed for use by agents that run on small-context models. Whereas an agent using a larger model may be able to keep in context numerous tools, smaller models do have this luxury, and must therefore be more parsimonious with the tools that they retain in context.

The solution that Laridae provides for this problem is enable registration, search, and retrieval of tools. A model may instead identify the general task it wants to perform, and then search for tools that perform that task. Laridae will then return an OpenAPI PathItem and server spec describing how to invoke that tool. After this point, the agent may eject any history regarding the tool search from its context.

Laridae also enables the set of tools available to agents to change over time. Clients may register new tools by causing Laridae to ingest their OpenAPI specs.

The original design statement(s) for Laridae may be [found here](https://www.leozqin.me/posts/a-registry-and-proxy-server-for-agentic-tool-calling/)

## Roadmap
The general roadmap for Laridae is certainly open to user feedback, but is broadly as follows:
- CRUD for tools
- Proxy server capability (to avoid shipping creds to your agents)
- Enhanced search and custom tagging

## To Run
Laridae is a Go application. It requires Go 1.23 or higher. Install dependencies by doing `go mod download`, then start the server by doing `go run .`

Read the docs at `http://localhost:8888/docs` to familiarize yourself with how to use Laridae.

Or, if you're feeling brave, pass the OpenAPI spec for Laridae straight to your agent and let them figure it out (`http://localhost:8888/openapi.json`)

This is really basic, and will surely evolve over time.
