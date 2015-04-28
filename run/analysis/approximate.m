function approximate
  use('Interaction');

  filename = locate('approximate');
  solution = h5read(filename, '/solution');
  surrogate = solution.Surrogate;

  ni = double(surrogate.Inputs);
  no = double(surrogate.Outputs);
  indices = reshape(surrogate.Indices{1}, ni, [])';
  surpluses = reshape(surrogate.Surpluses{1}, no, [])';

  Plot.figure(1200, 400);
  for i = 1:no
    subplot(no, 1, i);
    semilogy(abs(surpluses(:, i)));
    Plot.title('Output %d', i-1);
  end

  nn = size(indices, 1);

  nu = size(unique(indices, 'rows'), 1);
  if nu ~= nn
    warning('found %d nonunique indices out of %d', nn-nu, nn);
  end

  na = sum(admit(indices));
  if na ~= nn
    warning('found %d inadmissible indices out of %d', nn-na, nn);
  end
end
