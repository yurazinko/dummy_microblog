class Home extends React.Component {
  render() {
    return (
      <div className="container">
        <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
          <h1>Dummy Microblog</h1>
          <p>Sign in to get access </p>
          <a onClick={this.authenticate} className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
        </div>
      </div>
    )
  }
}

class LoggedIn extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      entries: []
    }
  }

  render() {
    return (
      <div className="container">
        <div className="col-lg-12">
          <br />
          <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
          <h2>Dummy microblog</h2>
          <p>Lets write something</p>
          <div className="row">
            {this.state.entries.map(function(entry, i){
              return (<Entry key={i} entry={entry} />);
            })}
          </div>
        </div>
      </div>
    )
  }
}

class Entry extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      liked: ""
    }
    this.like = this.like.bind(this);
  }

  like() {
    // ... we'll add this block later
  }

  render() {
    return (
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading">#{this.props.entry.id} <span className="pull-right">{this.state.liked}</span></div>
          <div className="panel-body">
            {this.props.entry.entry}
          </div>
          <div className="panel-footer">
            {this.props.entry.likes} Likes &nbsp;
            <a onClick={this.like} className="btn btn-default">
              <span className="glyphicon glyphicon-thumbs-up"></span>
            </a>
          </div>
        </div>
      </div>
    )
  }
}

class App extends React.Component {
  render() {
    if (this.loggedIn) {
      return (<LoggedIn />);
    } else {
      return (<Home />);
    }
  }
}

ReactDOM.render(<App />, document.getElementById("app"));
