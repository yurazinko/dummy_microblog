class App extends React.Component {
  
  render() {
    if (this.loggedIn) {
      return (<LoggedIn />);
    } else {
      return (<Home />);
    }
  }
}
