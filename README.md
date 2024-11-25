## Dependencies

- nodejs and npm
	- download nodejs from the [website](https://nodejs.org/en/download/prebuilt-binaries)
	- extract the file<br>
	`sudo tar -C /usr/local -xJf ~/Downloads/node-v22.11.0-linux-x64.tar.xz`
	- add the binaries to your path<br>
	`echo "PATH=/usr/local/bin/node-v22.11.0-linux-x64/bin:$PATH" >> ~/.profile`
	`source ~/.profile`

- typescript
```sh
# install typescripts locally in the project
npm install typescript
```

## Project Structure

The project's directory structure is organized as follows:

- **`conception/`**: Contains documents related to the architectural
		and design of the project.

- **`internal/`**: This is where the core logic of the application resides.
	It contains various packages:

	- **`config/`**: Contains every configurations functions
		and global constants.

	- **`database/`**: Contains every logics related to
		the database manipulation.

	- **`server/`**: Contains all server-related logics.

		- **`handlers/`**: Contains functions for handling HTTP requests.
			Each handler corresponds to an endpoint.

		- **`middleware/`**: Contains middleware functions that modify the
			request-response cycle, such as session management or logging.

		- **`models/`**: Contains structures related to the server.

		- **`services/`**: Contains functionality related to the handlers
			or specific to the business logic.

		- **`templates/`**: Contains basic functions related to the templates.
			Can also include specific functions that can be called inside of
			the template.

	- **`utils/`**: Utility functions that are used throughout the application.

- **`web/`**: Contains all the web-related files.

	- **`static/`**: Every files directly accessible by the browser.

		- **`scripts/`**: Additionnal scripts JS.

		- **`style/`**: Every CSS files.

	- **`templates/`**: Contains every html templates.

		- **`components`**: Other reusable elements.

		- **`layout/`**: Templates for the main layout (header, footer ...)

		- **`pages`**: Templates of the pages. Only those templates are called.

		- **`partials`**: Templates shared between several pages.
