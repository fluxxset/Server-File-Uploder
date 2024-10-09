Here is a proper `README.md` with Markdown formatting that you can copy and paste:

```markdown
# Server-File-Uploader

**Server-File-Uploader** is a simple tool designed to make file uploads to a server easier for users who may have difficulty using command-line methods such as SFTP or SCP. With this tool, you just run the binary, and it provides an easy-to-use web interface where you can upload files directly to the same folder where the binary is located.

## Features

- Simple and easy-to-use HTML interface for uploading files.
- No need for SFTP, SCP, or other command-line tools to upload files.
- Files are saved in the same folder as the binary (in a `web-uploads` directory created automatically).
- Cross-platform: Works wherever Go binaries can run (Linux, Mac, Windows).

## How to Use

1. **Build the application**:
   If you haven't built the binary yet, run the following command in the project directory:
   ```bash
   go build
   ```

2. **Run the server**:
   After building the binary, run it:
   ```bash
   ./Server-File-Uploader
   ```

3. **Upload your file**:
   - Open your web browser and navigate to `http://localhost:8080`.
   - Use the file upload form to upload your file.
   - The file will be saved in the `web-uploads` folder, created in the same directory where the binary is located.

4. **Stop the server**:
   Once you've uploaded your file, you can stop the server by pressing `Ctrl+C` in the terminal where the server is running.

---

## Important Warning

### ⚠️ **This tool uses a simple HTTP connection without any authentication**. For security reasons:

- Only run the server when you need to upload a file.
- **After you are done uploading, stop the server immediately by pressing `Ctrl+C` in the terminal**.
- Do not expose this tool on a public network without implementing additional security measures like HTTPS and authentication.

---

## License

This project is open-source under the MIT License.
```

You can now copy and paste this Markdown into your `README.md` file. Let me know if you need further changes!