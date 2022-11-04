class Footer extends React.Component {
    render() {
        if(!this.props.message){
            //Keine Meldung --> Keine Ausgabe
            return (<div/>);
        }
        return (
            <footer>
                <div className="row">
                    
                        <div className="col-2" />
                        <div className="col-4">Status</div>
                        <div className="col-4">{this.props.message}</div>

                    </div>
            </footer>
        );

    }
}