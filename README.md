## Dependencies

### nodejs and npm
1. First download nodejs from the [ official website](https://nodejs.org/en/download/prebuilt-binaries).

2. Install nodejs
```sh
rm -rf ~/.local/share/Trash/nodejs
[ -d /usr/local/nodejs ] && sudo mv /usr/local/nodejs ~/.local/share/Trash
sudo tar -C /usr/local -xJf ~/Downloads/node-v*-linux-x64.tar.xz
# MAKE SUR THAT THE LOCATION IS RIGHT ^^^^^^^^^^^^^^^^^^^^^^^^^^
sudo mv /usr/local/node-v*-linux-x64 /usr/local/nodejs
```

3. Add the binaries in your path.
```sh
sudo nano /etc/environment
```

4. Add ":/usr/local/nodejs/bin" at the end of the path.

5. And then update your path.
```sh
source /etc/environment
```


### typescript
```sh
npm install typescript
```

### sass
```sh
npm install sass
```

### openssl
- linux
```sh
sudo apt install openssl
```
- windows
1. [Download openssl](https://slproweb.com/products/Win32OpenSSL.html)
2. Include the bin folder into your path.


### gcc for Windows
1. [Download tdm-gcc](https://jmeubank.github.io/tdm-gcc/download/)
2. Include the bin folder into your path.

### make for Windows


## Run the application
```sh
make
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

		- **`services/`**: Contains functionality used in several handlers.

		- **`templates/`**: Contains the core logic of the template rendering and also
			include specific functions that can be called inside of the templates.

	- **`utils/`**: Utility functions that are used throughout the application.

- **`web/`**: Contains all the web-related files.

	- **`src/`**: Typescript files.

	- **`static/`**: Every files directly accessible by the browser.

		- **`scripts/`**: JS files generated from src/.

		- **`style/`**: Every CSS files.

	- **`templates/`**: Contains every html templates.

		- **`components/`**: Template with single reusable elements.

		- **`layout/`**: Templates for the main layout (header, footer ...)

		- **`pages/`**: Templates of the pages. Only those templates are called.

		- **`partials/`**: Templates shared between several pages.
