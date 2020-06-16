<p align="center">
  <a href="https://monitoror.com"> 
    <img src=".assets/monitoror-logo-and-text.svg" alt="" width="70%">
  </a>
</p>

<p align="center">
  <a href="https://github.com/monitoror/monitoror/actions?query=workflow%3Acontinuous-integration"><img src="https://img.shields.io/github/workflow/status/monitoror/monitoror/continuous-integration?style=for-the-badge" alt="Build"/></a>
  <a href="https://codecov.io/gh/monitoror/monitoror"><img src="https://img.shields.io/codecov/c/gh/monitoror/monitoror/master.svg?style=for-the-badge" alt="Code Coverage"/></a>
  <a href="https://github.com/monitoror/monitoror/releases"><img src="https://img.shields.io/github/release/monitoror/monitoror.svg?style=for-the-badge" alt="Releases"/></a>
  <br>
  <img src="https://img.shields.io/badge/Go-1.14-blue.svg?style=for-the-badge" alt="Go"/>
  <img src="https://img.shields.io/badge/NodeJS-10.0-blue.svg?style=for-the-badge" alt="NodeJS"/>
</p>

-----

# Monitoror

Monitoror is a wallboard monitoring app to monitor server status; monitor CI builds progress or even display critical values.


## Demo

<p align="center">
  <a href="https://demo.monitoror.com">
    <img src=".assets/monitoror-mockup.svg" alt="" width="70%"/> <br>
    Visit Monitoror live demo
  </a>
</p>


## Getting started

Monitoror is a single file app written in Go which can be run on these platforms:

- Linux (64bits, ARM)
- macOS
- Windows (64bits)

The app is divided into two parts: Core and UI.

Core is the server-side Monitoror HTTP API, configured by the environment variables or `.env` file.

UI is the client-side Monitoror loaded in browser, which is the wallboard itself.

[Visit the Monitoror website for more details](https://monitoror.com)


## Documentation

All details about [**installation**](https://monitoror.com/documentation/#installation) and [**configuration**](https://monitoror.com/documentation/#configuration) are on [our documentation](https://monitoror.com/documentation/)


## Development

See our [development guide](https://monitoror.com/guides/#development)


## Authors

<table>
<tbody>
  <tr width="100%">
    <td align="center" width="50%">
      <a href="https://github.com/jsdidierlaurent">
        <img src="https://avatars2.githubusercontent.com/u/11354381?s=150&v=4"><br>
        @jsdidierlaurent
      </a> <br>
      <strong>Jean-Sébastien Didierlaurent</strong><br>
      <em>Mostly on Monitoror Core</em><br>
      &bull; &bull; &bull;<br>
      https://twitter.com/Akhiro
    </td>
    <td align="center" width="50%">
      <a href="https://github.com/Alex-D">
        <img src="https://avatars2.githubusercontent.com/u/426843?s=150&v=4"><br>
        @Alex-D
      </a> <br>
      <strong>Alexandre Demode</strong><br>
      <em>Mostly on Monitoror UI</em><br>
      &bull; &bull; &bull;<br>
      https://twitter.com/AlexandreDemode
    </td>
  </tr>
</tbody>
</table>


## Support us

You can support Monitoror ongoing development by doing a donation or being a backer or a sponsor:

<a href="https://opencollective.com/monitoror/donate" target="_blank">
  <img src="https://opencollective.com/monitoror/donate/button@2x.png?color=blue" width="300" alt="Donate via Open Collective"/>
</a>


## License

This project is licensed under the MIT License - see [the LICENSE file](LICENSE) for details.
