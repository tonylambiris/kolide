var Kolide = {

  Request: {},

  table: false,
  fixedTable: false,

  selectedNodes: {},

  updateSelectedStatus: function() {
    var self = this;
    var html = "";

    if ($("#header .title").length > 0) {

      for (item in Kolide.selectedNodes) {

        if (Kolide.selectedNodes[item]) {
          var data = {
            name: item,
            type: "node"
          };

          html += Kolide.Templates.querySelectHeader(data);
        }
      }

      if (html.length <= 0) {
        html = "Query All Nodes";
      } else {
        html = "Query: " + html;
      }

      $("#header .title").html(html);
    }
  },

  FormatQuery: function(query) {
    query = query.replace(/\-\-.*$/gm, '');
    query = query.replace(/(\r\n|\n|\r)/gm, " ")

    split = query.split(";");

    if (split.length > 1) {

      var clean = []

      for (var i = 0, len = split.length; i < len; i++) {
        var item = split[i].trim();
        if (item !== "") {
          clean.push(item);
        }
      }

      query = clean.pop();
    }

    return query
  },

  lbox: {
    open: function(template, data, args) {
      var self = this;

      if (args && args.hasOwnProperty('width')) {
        if (args.width) {
          Kolide.lbox.options.style.width = args.width;
        }
      }

      var options = $.extend({}, Kolide.lbox.options, args);
      options.template = template;
      options.templateData = data;

      var box = $.limp(options);

      return box;
    },
    options: {
      cache: false,
      adjustmentSize: 0,
      loading: true,
      alwaysCenter: true,
      animation: "pop",
      shadow: "none",
      round: 0,
      distance: 10,
      overlayClick: true,
      enableEscapeButton: true,
      dataType: 'html',
      centerOnResize: true,
      closeButton: false,
      style: {
        '-webkit-outline': 0,
        color: '#000',
        position: 'fixed',
        border: '1px solid #ededed',
        outline: 0,
        zIndex: 10000001,
        opacity: 0,
        // overflow: 'auto',
        background: '#fff'
      },
      inside: {
        background: '#fff',
        padding: '0',
        display: 'block',
        border: 'none',
        overflow: 'visible'
      },
      overlay: {
        background: '#000',
        opacity: 0.7
      },
      onOpen: function() {},
      afterOpen: function() {},
      afterDestroy: function() {},
      onTemplate: function(template, data, limp) {
        return template;
      }
    }
  },

  Node: {
    current: null,
    tables: [],

    Ask: function(id, question, callback) {
      var json = {}

      Kolide.Socket.request(question, id, function(err, data) {
        if (typeof callback === "function") {
          return callback(data, err);
        }
      });
    },

    commands: {
      confirm: function(data, callback) {
        var box = Kolide.lbox.open(Kolide.Templates.confirm(data), {}, {
          width: "760px",
          disableDefaultAction: true,
          afterClose: function() {},
          afterOpen: function() {},
          onAction: function() {
            callback();
          }
        });

        box.open();
      },
      disconnect: function(data) {
        this.confirm(data, function() {
          Kolide.Socket.request('disconnect', data.id, function(err) {
            if (err) {
              Kolide.Flash.error(err);
            }

            $.limpClose();
          });
        });
      },
      delete: function(data) {
        this.confirm(data, function() {
          $.ajax({
            url: "/api/v1/nodes/" + data.id,
            type: "delete",
            success: function(d) {
              if (d.err) {
                Kolide.Flash.error(d.err);
              } else {
                $("li.node[data-node-id='" + d.context.node.key + "']").remove();
              }

              $.limpClose();
            }, error: function() {
              Kolide.Flash.error(err);
            }
          })
        });
      }
    },

    fetchTables: function(callback) {
      var self = this;

      Kolide.Socket.request('tables', {
        id: self.current
      }, function(err, data) {

        self.tables = data;

        if (typeof callback === "function") {
          return callback(data, err);
        }
      });
    },

    fetchTableInfo: function(table, callback) {
      var self = this;

      if (Kolide.fixedTable) {
        Kolide.fixedTable._fnDestroy();
        Kolide.fixedTable = false;
      }

      if (Kolide.table) {
        Kolide.table.destroy();
        Kolide.table = false;
      }

      Kolide.Loading.start();

      Kolide.Socket.request('table-info', {
        id: self.current,
        sql: "pragma table_info(" + table + ");",
      }, function(err, data) {

        Kolide.Editor.self.setValue("select * from " + table + ";");

        if (typeof callback === "function") {
          data.hideNode = true;
          Kolide.Query.Render([data], err, function() {
            return callback(data, err);
          });
        }

      });
    },

    clean: function() {
      this.current = null;
      this.tables = [];

      if (Kolide.fixedTable) {
        Kolide.fixedTable._fnDestroy();
        Kolide.fixedTable = false;
      }

      if (Kolide.table) {
        Kolide.table.destroy();
        $(".query-results").remove();
        Kolide.table = false;
      }

      $("#content").removeClass("node-view");
      $("#node-tables").remove();
      // $(".save-query, .load-query").addClass("disabled");
      $('.table-filter').addClass("node-view");
      $(".table-filter").show().find("input").val("");
      $("#content").css("left", 460);
    },

    close: function() {
      this.clean();
      $("#header .title").text("Query All Nodes");
      $("a.run-query").text("Run Query");
      // $(".save-query, .load-query").removeClass("disabled");
      $('.table-filter').removeClass("node-view");
      $(".table-filter").hide().find("input").val("");
      $("#content").css("left", 230);
    },

    open: function(name, id) {
      var self = this;

      this.clean();
      this.current = id;

      Kolide.Loading.start();
      this.fetchTables(function(data, err) {
        Kolide.Loading.done();

        if (err) {
          Kolide.Flash.error(err);

          Kolide.Socket.request('disconnect-dead', id, function(err) {
            if (err) {
              Kolide.Flash.error(err);
            }
          });
          return;
        }

        $("a.run-query").text("Query Node");
        $("#header .title").text("Query Node: " + data.name + " (" + data.hostname + ")");
        $("#content").addClass("node-view");
        $("#wrapper").append(Kolide.Templates.tables(data.results));

        $("ul.tables li").on("click", function(e) {
          e.preventDefault();

          $("ul.tables li").removeClass("selected");
          $(this).addClass("selected");

          var table = $(this).attr("data-table-name");
          self.fetchTableInfo(table, function() {})
        });

        $("ul.tables li:first-child").click();
      });

    }
  },

  Flash: {
    delay: 1500,

    show: function(data, type) {
      var self = this;

      if (self.timeout) {
        clearTimeout(self.timeout);
        self.timeout = null
        self.hide();
      }

      self.hide();

      $("#flash-message").attr("class", type).text(data).show();

      self.timeout = setTimeout(function() {
        $("#flash-message").stop().fadeOut("slow", function() {
          self.hide();
        });
      }, this.delay);
    },
    hide: function() {
      $("#flash-message").attr("class", "").text("").hide();
    },
    error: function(message) {
      this.show(message, "error");
    },
    success: function(message) {
      this.show(message, "success");
    }
  },

  Loading: {
    options: {
      catchupTime: 100,
      initialRate: .03,
      minTime: 250,
      ghostTime: 100,
      maxProgressPerFrame: 20,
      easeFactor: 1.25,
      startOnPageLoad: true,
      restartOnPushState: true,
      restartOnRequestAfter: 500,
      target: 'body',
      elements: {
        checkInterval: 100,
        selectors: ['body']
      },
      eventLag: {
        minSamples: 10,
        sampleCount: 3,
        lagThreshold: 3
      },
      ajax: {
        trackMethods: ['GET', 'POST', 'DELETE'],
        trackWebSockets: true,
        ignoreURLs: []
      }
    },
    start: function() {
      this.self = Pace.start(this.options);
      $("#content .wrapper").css("opacity", 0.5);
      $("#content .loading").show();
      // $("#loading").show();
    },
    stop: function() {
      Pace.stop();
      $("#content .wrapper").css("opacity", 1);
      $("#content .loading").hide();
      // $("#loading").hide();
    },
    restart: function() {
      Pace.restart();
    },
    done: function() {
      Pace.stop();
      $("#content .wrapper").css("opacity", 1);
      $("#content .loading").hide();
      // $("#loading").hide();
    }
  },

  Templates: {

    Init: function() {
      this.table = Handlebars.compile($("#query-results-table").html());
      this.row = Handlebars.compile($("#query-results-row").html());
      this.node = Handlebars.compile($("#node-template").html());
      this.tables = Handlebars.compile($("#tables-template").html());
      this.saveQuery = Handlebars.compile($("#save-query-template").html());
      this.loadQuery = Handlebars.compile($("#load-query-template").html());
      this.nodeContextMenu = Handlebars.compile($("#node-context-menu-template").html());
      this.confirm = Handlebars.compile($("#confirm-template").html());

      this.querySelectHeader = Handlebars.compile($("#query-selected-header").html());

      Handlebars.registerHelper("formatDate", function(datetime, format) {
        if (moment) {
          // can use other formats like 'lll' too
          return moment(datetime).format(format);
        }
        else {
          return datetime;
        }
      });

      Handlebars.registerHelper ('truncate', function(str, len) {
        if (str.length > len) {
          var new_str = str.substr (0, len + 1);

          while (new_str.length) {
            var ch = new_str.substr (-1);
            new_str = new_str.substr (0, -1);

            if (ch == ' ') {
              break;
            }
          }

          if (new_str == '') {
            new_str = str.substr (0, len);
          }

          return new Handlebars.SafeString (new_str + '...');
        }
        return str;
      });
    }

  },

  Query: {

    RunSavedQuery: function(name, query, callback) {
      $("#content").scrollTop(0);

      $("#header .title").text("Saved Query: " + name);
      Kolide.Editor.self.setValue(query);

      if (Kolide.fixedTable) {
        Kolide.fixedTable._fnDestroy();
        Kolide.fixedTable = false;
      }

      if (Kolide.table) {
        Kolide.table.destroy();
        Kolide.table = false;
      }

      Kolide.Loading.start()

      Kolide.Query.Run("query", Kolide.FormatQuery(query), function(results, err) {
        Kolide.Query.Render(results, err);
        if (typeof callback === "function") {
          callback();
        }
      });
    },

    Save: function(params) {
      var editor;

      params.query = Kolide.Editor.self.getValue();

      var box = Kolide.lbox.open(Kolide.Templates.saveQuery(params), {}, {
        width: "760px",
        disableDefaultAction: true,
        afterClose: function() {
          editor.destroy();
        },
        afterOpen: function() {
          editor = Kolide.Editor._build("save-query-editor");
          editor.setValue(params.query);
          $("#save-query-name").focus()
        },
        onAction: function() {

          if (Kolide.Request.save) {
            Kolide.Request.save.abort()
          }

          var saveParams = {
            name: $("#save-query-name").val(),
            query: editor.getValue(),
            type: "all"
          }

          if (saveParams.name.length <= 0) {
            $("#save-query-name").addClass("error");
            return
          } else {
            $("#save-query-name").removeClass("error");
          }

          if (saveParams.query.length <= 0) {
            $("#save-query-editor").addClass("error");
            return;
          } else {
            $("#save-query-editor").removeClass("error");
          }

          Kolide.Request.save = $.ajax({
            url: "/api/v1/saved-queries",
            dataType: 'json',
            contentType: "application/json",
            type: "POST",
            data: JSON.stringify(saveParams),
            success: function(data) {
              $.limpClose();

              if (data.error && data.error.length > 0) {
                Kolide.Flash.error(data.error);
              } else {
                Kolide.Flash.success("Query saved successfully.");
              }

            },
            error: function(a, b, c) {
              // console.log(a,b,c)
            }
          })
        }
      });

      box.open();
    },

    Load: function(params) {

      if (Kolide.Request.load) {
        Kolide.Request.load.abort();
      }

      Kolide.Request.load = $.ajax({
        url: "/api/v1/saved-queries",
        dataType: "json",
        contentType: "multipart/form-data",
        error: function(a, b, c) {
          // console.log(a,b,c)
        },
        success: function(data) {
          Kolide.Request.load = null;

          var box = Kolide.lbox.open(Kolide.Templates.loadQuery(data.context.queries), {}, {
            width: "800px",
            disableDefaultAction: true,
            afterClose: function() {},
            afterOpen: function() {
              var queryList = new List('load-query-select', {
                valueNames: ['name', 'query']
              });

              $("a.load-saved-query").on("click", function(e) {
                e.preventDefault();
                var item = $(this).parents("li");
                var query = item.find(".query").text();
                var name = item.find(".name").text();

                Kolide.Query.RunSavedQuery(name, query);
                $.limpClose();
              });

              $("a.delete-saved-query").on("click", function(e) {
                e.preventDefault();

                var item = $(this).parents("li");
                var id = item.attr("data-query-id");

                $.ajax({
                  url: "/api/v1/saved-queries/" + id,
                  type: "DELETE",
                  dataType: "json",
                  data: {
                    id: id
                  },
                  success: function(data) {
                    item.fadeOut("fast");
                  }
                });
              });
            },
            onAction: function() {}
          });

          box.open();
        }
      });

    },

    Execute: function() {
      $("#content").scrollTop(0);

      if (Kolide.fixedTable) {
        Kolide.fixedTable._fnDestroy();
        Kolide.fixedTable = false;
      }

      if (Kolide.table) {
        Kolide.table.destroy();
        Kolide.table = false;
      }

      var value = Kolide.Editor.self.getValue();

      Kolide.Loading.start()

      Kolide.Query.Run("query", Kolide.FormatQuery(value), function(results, err) {
        Kolide.Query.Render(results, err);
      });

    },

    Render: function(results, err, callback) {

      var err = false;
      var table = null;
      var count = 0

      if (!results) {
        Kolide.Flash.error("Your query returned no data.")
        return
      }

      for (node in results.context.results) {
        var n = results.context.results[node];

        if (n.timeout || !n.results) {
          continue;
        }

        if (!table) {
          var data = {
            hideNode: false,
            name: n.node.name,
            hostname: n.node.address,
            results: n.results[0]
          }

          table = Kolide.Templates.table(data);
          $("#content .wrapper").html(table);
        };


        for (row in n.results) {
          var r = n.results[row];

          var data = {
            hideNode: false,
            name: n.node.name,
            hostname: n.node.address,
            results: r
          }

          var row = Kolide.Templates.row(data)
          $("table.query-results tbody").append(row);

          count++;
        }
      }

      if (count <= 0) {
        if (Kolide.fixedTable) {
          Kolide.fixedTable._fnDestroy();
          Kolide.fixedTable = false;
        }

        if (Kolide.table) {
          Kolide.table.destroy();
          Kolide.table = false;
        }

        Kolide.Loading.done()
        Kolide.Flash.error("No results found.");

        $("#content .wrapper").html("");
        return
      }

      Kolide.table = $("table.query-results")
        .on('order.dt', function() {
          if (Kolide.fixedTable) {
            Kolide.fixedTable.fnUpdate()
            $("#content").scrollTop(0);
          }
        }).DataTable({
          searching: true,
          paging: false,
          info: false,
          ordering: true,
          order: []
        });

        Kolide.fixedTable = new $.fn.dataTable.FixedHeader(Kolide.table, {
        });

        $(".table-filter").show().find("input").val("");

        window.onresize = function() {
          if (Kolide.fixedTable) {
            Kolide.fixedTable.fnUpdate()
          }
        }


        Kolide.Editor.self.focus();
        Kolide.Loading.done()


        if (typeof callback === "function") {

          return callback();
        }
    },

    Run: function(type, sql, callback) {
      var all = true;

      var nodes = [];

      for (node in Kolide.selectedNodes) {
        if (Kolide.selectedNodes[node]) {
          nodes.push(node);
        }
      }

      if (nodes.length > 0) {
        all = false;
      }

      $.ajax({
        url: "/api/v1/query",
        type: "post",
        dataType: 'json',
        contentType: "application/json",
        timeout: 30 * 1000,
        data: JSON.stringify({
          nodes: nodes,
          all: all,
          sql: sql
        }),
        success: function(data) {
          Kolide.Loading.done();

          if (typeof callback === "function") {
            return callback(data, null);
          }

        }, error: function(a, b, c) {
          console.log(a,b,c);
          return callback(null, b);
        }
      })

    }
  },

  Editor: {
    self: null,

    _build: function(div) {
      ace.require("ace/ext/language_tools");

      var editor = ace.edit(div);

      editor.setOptions({
        enableBasicAutocompletion: true,
        enableSnippets: false,
        enableLiveAutocompletion: true
      });

      editor.getSession().setMode("ace/mode/sql");
      editor.getSession().setTabSize(2);
      editor.getSession().setUseSoftTabs(true);
      editor.getSession().setUseWrapMode(true);
      editor.setHighlightActiveLine(false);
      editor.setShowPrintMargin(false);
      // document.getElementById('editor').style.fontSize='13px';

      return editor;
    },

    Build: function() {

      this.self = this._build("editor");

      this.self.focus();
      // this.self.setValue("select * from listening_ports a join processes b on a.pid = b.pid;");
      this.self.setValue("select * from listening_ports join processes using (pid);");

      this.self.commands.addCommands([
        {
          name: "run_query",
          bindKey: {
            win: "Ctrl-Enter",
            mac: "Command-Enter"
          },
          exec: function(editor) {
            Kolide.Query.Execute();
          }
        }
      ]);

      $("a.run-query").on("click", function(e) {
        e.preventDefault();
        Kolide.Query.Execute();
      });

      $("a.export-results").on("click", function(e) {
        e.preventDefault();
        var csv = $("table.query-results").table2CSV({
          delivery: 'value'
        });
        window.location.href = 'data:text/csv;charset=UTF-8,'
          + encodeURIComponent(csv);
      });
    }
  },

  Socket: null,
  Init: function() {
    // fetch the CSRF token from the meta tag
    var token = $("meta[name='_csrf']").attr("content");

    // ensure every Ajax request has the CSRF token
    // included in the request's header.
    $(document).ajaxSend(function(e, xhr, options) {
      xhr.setRequestHeader("X-CSRF-TOKEN", token);
    });

    self.Socket = new WebSocket("wss://"+window.location.host+"/api/v1/websocket");

    var updateNode = function(node) {
      var item = $("li[data-node-id='" + node.data.key + "']");

      var selected = false;

      if (Kolide.selectedNodes.hasOwnProperty(node.data.key)) {
        if (Kolide.selectedNodes[node.data.key])  {
          selected = true;
        }
      }

      node.data.selected = selected;

      if (item.length > 0) {
        item.replaceWith(Kolide.Templates.node(node.data))
      } else {
        $("ul#nodes").append(Kolide.Templates.node(node.data))
      }

      delete (Kolide.nodeList)

      Kolide.nodeList = new List('sidebar', {
        valueNames: [
          'node-name', 'node-ip', 'node-node-id', 'online-status'
        ]
      });
    }

    self.Socket.onclose = function(evt) {}

    self.Socket.onmessage = function(evt) {
      var json = JSON.parse(evt.data);

      if (json.type === "node") {
        updateNode(json);
      }
    }

    // if (!!window.EventSource) {
    // var source = new EventSource('/api/v1/query');

    // source.addEventListener('message', function(e) {
    // console.log(e)
    // }, false);

    // } else {
    // alert("NOT SUPPORTED");
    // }

    Kolide.Templates.Init()
    Kolide.Editor.Build()
  }
};

jQuery(document).ready(function($) {
  Kolide.Init();

  $.fn.dataTableExt.sErrMode = "none";

  var lastScrollLeft = 0;
  $("#content").on("scroll", function() {
    if (Kolide.fixedTable) {
      var documentScrollLeft = $("#content").scrollLeft();
      if (lastScrollLeft != documentScrollLeft) {
        // super hack
        Kolide.fixedTable.fnPosition();
        $(".FixedHeader_Cloned").css("top", 230);
        lastScrollLeft = documentScrollLeft;
      }
    }
  });

  $(document).on("click", "li.node", function(e) {
    e.preventDefault();

    var name = $(this).find("span.node-name").text();
    var id = $(this).attr("data-node-id");

    if ($(this).hasClass("online")) {
      if ($(this).hasClass("current")) {

        $("li.node").removeClass("current");
        Kolide.selectedNodes[id] = false;
      } else {

        $("li.node").removeClass("current");
        $(this).addClass("current");
        Kolide.selectedNodes[id] = true;
      }

    } else {
      Kolide.Flash.error("Node (" + name + ") is current offline.");
    }

    Kolide.updateSelectedStatus();
  });

  $(document).on("contextmenu", "li.node", function(event) {
    event.preventDefault();

    $("ul.custom-menu").remove();

    var parent = $(this);
    var id = parent.attr("data-node-id");

    var menu = $(Kolide.Templates.nodeContextMenu({
      id: id
    })).appendTo("body")
    .css({
      top: event.pageY + "px",
      left: event.pageX + "px"
    });

    $("a.disconnect-node").on("click", function(e) {
      e.preventDefault();

      if (parent.hasClass("online")) {
        Kolide.Node.commands.disconnect({
          id: id,
          title: "Disconnect Node",
          message: "Are you sure you want to disconnect this node?",
          icon: "fa-exclamation",
          buttonTitle: "Disconnect"
        });
      } else {
        Kolide.Flash.error("This node is already offline.");
      }
    });

    $("a.system-information-node").on("click", function(e) {
      e.preventDefault();

      if (parent.hasClass("online")) {

        Kolide.Loading.start();

        Kolide.Node.Ask(id, "system-information", function(data, err) {

          Kolide.Loading.done();

          if (err) {
            Kolide.Flash.error(err);
            return
          }

          var box = Kolide.lbox.open(Kolide.Templates.systemInformation(data), {}, {
            disableDefaultAction: true,
            afterClose: function() {},
            afterOpen: function() {},
            onAction: function() {
            }
          });

          box.open();
        });
      } else {
        Kolide.Flash.error("This node is offline.");
      }
    });

    $("a.delete-node").on("click", function(e) {
      e.preventDefault();
      Kolide.Node.commands.delete({
        id: id,
        title: "Delete Node",
        message: "Are you sure you want to delete this node?",
        icon: "fa-trash",
        buttonTitle: "Delete"
      });
    });

  }).on("click", function(event) {
    $("ul.custom-menu").hide();
  });

  Kolide.nodeList = new List('sidebar', {
    valueNames: [
      'node-name', 'node-ip', 'node-node-id', 'online-status'
    ]
  });

  Kolide.nodeList.sort('online-status', {
    order: "desc"
  });

  $(document).on("click", ".envdb-control a.save-query", function(e) {
    e.preventDefault();
    Kolide.Query.Save({});
  });

  $(document).on("click", ".envdb-control a.load-query", function(e) {
    e.preventDefault();
    Kolide.Query.Load({});
  });

  $(document).on("input", ".table-filter input", function(e) {
    var value = $(this).val();
    // search
    if (Kolide.table) {
      Kolide.table.search(value);
      Kolide.table.draw();
      Kolide.fixedTable.fnUpdate();
    }
  });

  $(document).on("click", "a.logout", function(e) {
    e.preventDefault();

    $.ajax({
      url: "/authorize",
      type: "DELETE",
      success: function() {
        window.location = "/";
      }, error: function(a, b, c) {
        Kolide.Flash.error("Error - " + b);
        console.log(a, b, c);
      }
    })
  });

  $(document).on("click", "i.hide-sidebar", function(e) {
    e.preventDefault();
    $("#sidebar").css("left", -230)
    $("#header, #envdb-query, #node-tables").css("left", 0);

    if ($("#content").hasClass("node-view")) {
      $("#content").css("left", 230);
    } else {
      $("#content").css("left", 0);
    }

    $("div.show-sidebar").show();

    if (Kolide.table) {
      Kolide.table.draw();

      if (Kolide.fixedTable) {
        Kolide.fixedTable.fnUpdate()
        Kolide.fixedTable.fnPosition();
      }
    }
  });

  $(document).on("click", "div.show-sidebar", function(e) {
    e.preventDefault();
    $("#sidebar").css("left", 0)
    $("#header, #envdb-query, #node-tables").css("left", 230);
    $("div.show-sidebar").hide();

    if ($("#content").hasClass("node-view")) {
      $("#content").css("left", 460);
    } else {
      $("#content").css("left", 230);
    }

    if (Kolide.table) {
      Kolide.table.draw();

      if (Kolide.fixedTable) {
        Kolide.fixedTable.fnUpdate()
        Kolide.fixedTable.fnPosition();
      }
    }
  });

});
