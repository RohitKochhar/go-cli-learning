# Chapter 3

## Overview
- In this chapter, we develop a tool to preview Markdown files locally, using a web browser
- The tool will convert Markdown source into HTMl that can be viewed in a browser by doing the following:
    1. Read the contents of the input MD file
    2. Use some Go external libraries to parse MD and generate valid HTML
    3. Wrap the results with an HTML header and footer
    4. Save the bugger to an HTML file that can be viewed in a browser
- This tool uses `blackfriday` to convert Markdown to HTML, and `bluemonday` to sanitize the output to ensure no malicious content
- In this chapter, we use a different code pattern to allow for more specific unit tests of each function, unlike in chapter 2 where the main function was testing directly.
- The testing approach used in this chapter is the _golden files_ technique, where the expected resutls are saved into files that are loaded during the tests for validating the actual output