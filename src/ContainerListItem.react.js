var $ = require('jquery');
var React = require('react/addons');
var Router = require('react-router');
var remote = require('remote');
var dialog = remote.require('dialog');
var ContainerStore = require('./ContainerStore');
var OverlayTrigger = require('react-bootstrap').OverlayTrigger;
var Tooltip = require('react-bootstrap').Tooltip;

var ContainerListItem = React.createClass({
  handleItemMouseEnter: function () {
    var $action = $(this.getDOMNode()).find('.action');
    $action.show();
  },
  handleItemMouseLeave: function () {
    var $action = $(this.getDOMNode()).find('.action');
    $action.hide();
  },
  handleDeleteContainer: function () {
    dialog.showMessageBox({
      message: 'Are you sure you want to delete this container?',
      buttons: ['Delete', 'Cancel']
    }, function (index) {
      if (index === 0) {
        ContainerStore.remove(this.props.container.Name, function (err) {
          console.error(err);
          var containers = ContainerStore.sorted();
          if (containers.length === 1) {
            $(document.body).find('.new-container-item').parent().fadeIn();
          }
        });
      }
    }.bind(this));
    return false;
  },
  render: function () {
    var self = this;
    var container = this.props.container;
    var imageNameTokens = container.Config.Image.split('/');
    var repo;
    if (imageNameTokens.length > 1) {
      repo = imageNameTokens[1];
    } else {
      repo = imageNameTokens[0];
    }
    var imageName = (
      <OverlayTrigger placement="bottom" overlay={<Tooltip>{container.Config.Image}</Tooltip>}>
        <div>{repo}</div>
      </OverlayTrigger>
    );

    // Synchronize all animations
    var style = {
      WebkitAnimationDelay: (self.props.start - Date.now()) + 'ms'
    };

    var state;
    if (container.State.Downloading) {
      state = (
        <OverlayTrigger placement="bottom" overlay={<Tooltip>Downloading</Tooltip>}>
          <div className="state state-downloading">
            <div style={style} className="downloading-arrow"></div>
          </div>
        </OverlayTrigger>
      );
    } else if (container.State.Running && !container.State.Paused) {
      state = (
        <OverlayTrigger placement="bottom" overlay={<Tooltip>Running</Tooltip>}>
          <div className="state state-running"><div style={style} className="runningwave"></div></div>
        </OverlayTrigger>
      );
    } else if (container.State.Restarting) {
      state = (
        <OverlayTrigger placement="bottom" overlay={<Tooltip>Restarting</Tooltip>}>
          <div className="state state-restarting" style={style}></div>
        </OverlayTrigger>
      );
    } else if (container.State.Paused) {
      state = (
        <OverlayTrigger placement="bottom" overlay={<Tooltip>Paused</Tooltip>}>
          <div className="state state-paused"></div>
        </OverlayTrigger>
      );
    } else if (container.State.ExitCode) {
      state = (
        <OverlayTrigger placement="bottom" overlay={<Tooltip>Stopped</Tooltip>}>
          <div className="state state-stopped"></div>
        </OverlayTrigger>
      );
    } else {
      state = (
        <OverlayTrigger placement="bottom" overlay={<Tooltip>Stopped</Tooltip>}>
          <div className="state state-stopped"></div>
        </OverlayTrigger>
      );
    }

    return (
      <Router.Link data-container={name} to="containerDetails" params={{name: container.Name}}>
        <li onMouseEnter={self.handleItemMouseEnter} onMouseLeave={self.handleItemMouseLeave}>
          {state}
          <div className="info">
            <div className="name">
              {container.Name}
            </div>
            <div className="image">
              {imageName}
            </div>
          </div>
          <div className="action">
            <span className="icon icon-delete-3 btn-delete" onClick={this.handleDeleteContainer}></span>
          </div>
        </li>
      </Router.Link>
    );
  }
});

module.exports = ContainerListItem;
