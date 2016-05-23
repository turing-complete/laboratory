function [filename, base] = locate(type)
  pattern = sprintf('_%s.h5', type);

  files = dir(pwd);
  for i = 1:length(files)
    filename = files(i).name;
    if ~isempty(strfind(filename, pattern))
      tokens = regexp(filename, '^(\d+_\d+_[a-z]+)_', 'tokens');
      base = tokens{1}{1};
      return
    end
  end

  error('Expected to find a file of type %s.', type);
end
