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

		- **`routes/`**: Contains every route endpoints of the application.

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

		- **`layout/`**: Contains templates that are common
			to several other templates like an header or a footer.
