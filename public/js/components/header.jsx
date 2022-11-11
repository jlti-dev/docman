class Header extends React.Component {
    render() {

        return (
            <header>
                <div className="row">
                    <div className="col-2">

                    </div>
                    <div className="col-2">
                        {this.props.app.state.role === "Admin" &&
                            <button onClick={() => this.props.app.setNavigation("upload")}>Upload</button>
                        }
                    </div>
                    <div className="col-2">
                        {this.props.app.state.role === "Admin" &&
                            <button onClick={() => this.props.app.setNavigation("user")}>Benutzer administrieren</button>
                        }
                    </div>
                    <div className="col-2">
                        <button onClick={() => this.props.app.setNavigation("download")}>Downloads</button>
                    </div>
                    <div className="col-2">
                        <button onClick={() => this.props.app.setNavigation("profile")}>Profil</button>
                    </div>
                    <div className="col-2">
                        <button onClick={() => this.props.app.setToken(undefined)}>Logout</button>
                    </div>
                </div>
            </header>
        );

    }
}