function filename = locale(type)
  pattern = sprintf('_%s.h5', type);

  files = dir(pwd);
  for i = 1:length(files)
    filename = files(i).name;
    if ~isempty(strfind(filename, pattern))
      return
    end
  end

  error('Expected to find a file of type %s.', type);
end
