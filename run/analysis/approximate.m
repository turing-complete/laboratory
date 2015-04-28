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

  good = admissibility(indices);
  if any(~good)
    warning('found %d inadmissible indices out of %d', sum(~good), length(good));
  end
end

function good = admissibility(indices)
  [nn, ni] = size(indices);

  good = false(nn, 1);
  indices = bitshift(bitshift(indices, 32), -32);
  indices = int32(indices);

  for i = 1:nn
    I = repmat(indices(i, :), ni, 1) - int32(eye(ni));
    I(any(I < 0, 1), :) = [];
    good(i) = all(ismember(I, indices, 'rows'));
  end
end
