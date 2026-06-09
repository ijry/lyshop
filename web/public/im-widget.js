(function () {
  var current = document.currentScript || {};
  var defaults = {
    baseUrl: current.getAttribute ? (current.getAttribute('data-base-url') || '') : '',
    token: current.getAttribute ? (current.getAttribute('data-token') || '') : '',
    position: 'right',
    title: '在线客服',
    launcherText: '客服',
    width: 420,
    height: 640,
    context: {}
  };

  function merge(target, source) {
    source = source || {};
    Object.keys(source).forEach(function (key) {
      if (source[key] !== undefined && source[key] !== null) target[key] = source[key];
    });
    return target;
  }

  function inferContext(extra) {
    return merge({
      visitor_id: localStorage.getItem('lyshop_im_visitor_id') || '',
      visitor_language: navigator.language || '',
      visitor_referrer: document.referrer || '',
      visitor_url: location.href,
      visitor_device: /Mobi|Android|iPhone|iPad/i.test(navigator.userAgent) ? 'mobile' : 'desktop',
      visitor_extra: {
        title: document.title || '',
        user_agent: navigator.userAgent || ''
      }
    }, extra || {});
  }

  function ensureVisitorID(context) {
    if (context.visitor_id) return context;
    var id = 'v_' + Date.now().toString(36) + '_' + Math.random().toString(36).slice(2, 10);
    localStorage.setItem('lyshop_im_visitor_id', id);
    context.visitor_id = id;
    return context;
  }

  function init(options) {
    var opts = merge(merge({}, defaults), options || {});
    var baseUrl = (opts.baseUrl || location.origin).replace(/\/$/, '');
    var context = ensureVisitorID(inferContext(opts.context));
    var root = document.createElement('div');
    var side = opts.position === 'left' ? 'left' : 'right';
    root.style.cssText = 'position:fixed;bottom:20px;' + side + ':20px;z-index:2147483000;font-family:ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;';

    var frame = document.createElement('iframe');
    var tokenParam = opts.token ? '&token=' + encodeURIComponent(opts.token) : '';
    frame.src = baseUrl + '/chat?embed=1' + tokenParam;
    frame.title = opts.title;
    frame.allow = 'clipboard-write';
    frame.style.cssText = 'display:none;width:' + Number(opts.width || 420) + 'px;height:' + Number(opts.height || 640) + 'px;max-width:calc(100vw - 40px);max-height:calc(100vh - 96px);border:0;border-radius:18px;box-shadow:0 24px 70px rgba(15,23,42,.24);background:#fff;overflow:hidden;';

    var button = document.createElement('button');
    button.type = 'button';
    button.textContent = opts.launcherText;
    button.style.cssText = 'margin-top:12px;float:' + side + ';border:0;border-radius:999px;background:#dc2626;color:#fff;padding:12px 18px;font-size:14px;font-weight:700;box-shadow:0 12px 34px rgba(220,38,38,.35);cursor:pointer;';

    button.onclick = function () {
      frame.style.display = frame.style.display === 'none' ? 'block' : 'none';
      if (frame.style.display === 'block') {
        frame.contentWindow && frame.contentWindow.postMessage({ type: 'lyshop-im-context', context: context }, baseUrl);
      }
    };
    frame.onload = function () {
      frame.contentWindow && frame.contentWindow.postMessage({ type: 'lyshop-im-context', context: context }, baseUrl);
    };

    root.appendChild(frame);
    root.appendChild(button);
    document.body.appendChild(root);
    return {
      open: function () { frame.style.display = 'block'; },
      close: function () { frame.style.display = 'none'; },
      destroy: function () { root.parentNode && root.parentNode.removeChild(root); }
    };
  }

  window.LYShopIMWidget = { init: init };
})();
