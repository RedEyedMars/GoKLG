function createTextLinks_(text) {
  return (text || "")
    .replace(/\{([^\{\}]+)\}/ig,  function(match, curl){ return "<"+curl+">"})
    .replace(/\{\{([^\}]+)\}\}/ig,function(match, curl){ return "{"+curl+"}"})
    .replace(/([^\S]|^)((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)/gi,
    function(match, space, url){
      var hyperlink = url;
      if (!hyperlink.match('^https?:\/\/')) {
        hyperlink = 'http://' + hyperlink;
      }
      if(hyperlink.match('https?:\/\/(www.)?youtu\.?be')){
        var match = /https?:\/\/(www.)?youtu(\.be|be\.com)\/(watch\?v=|embed\/)?([^&]+)(&[^&]+)*/g.exec(hyperlink);
        return space + '<iframe width="560" height="315" src="https://www.youtube.com/embed/'+match[4]+'" frameborder="0" allowfullscreen></iframe>';
      }
      else {
        return space + '<a href="' + hyperlink + '" target="_blank">' + url + '</a>';
      }
    });
};
