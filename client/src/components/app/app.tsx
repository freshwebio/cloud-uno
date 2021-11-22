import { Navbar, Nav, Container } from 'react-bootstrap'
import logo from '../../logo.svg'
import styles from './app.module.scss'

function App() {
  return (
    <div className={styles.App}>
      <header>
        <Navbar bg="primary" variant="dark">
          <Container>
            <Navbar.Brand href="#home">
              <img src={logo} alt="logo" width={150} />
            </Navbar.Brand>
            <Nav className="me-auto">
              <Nav.Link href="#home">Home</Nav.Link>
              <Nav.Link href="#features">Features</Nav.Link>
              <Nav.Link href="#pricing">Pricing</Nav.Link>
            </Nav>
          </Container>
        </Navbar>

        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className={styles['App-link']}
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  )
}

export default App
