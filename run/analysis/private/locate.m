function [files, names] = locate(type)
  pattern = sprintf('_%s.h5', type);
  entries = dir(pwd);
  files = {};
  names = {};
  for i = 1:length(entries)
    path = entries(i).name;
    if ~isempty(strfind(path, pattern))
      tokens = regexp(path, '^(\d+_\d+_[a-z]+)_', 'tokens');
      files{end + 1} = path;
      names{end + 1} = tokens{1}{1};
    end
  end
end
