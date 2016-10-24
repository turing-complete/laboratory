function [files, names] = locate(type, name)
  if nargin < 2
    pattern = sprintf('^(\\d+_\\d+_[a-z]+)_.*_%s.h5', type);
  else
    pattern = sprintf('^(%s)_.*_%s.h5', name, type);
  end
  entries = dir(pwd);
  files = {};
  names = {};
  for i = 1:length(entries)
    path = entries(i).name;
    takens = regexp(path, pattern, 'tokens');
    if ~isempty(takens)
      files{end + 1} = path;
      names{end + 1} = takens{1}{1};
    end
  end
end
